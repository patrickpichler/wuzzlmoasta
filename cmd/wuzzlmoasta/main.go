package main

//go:generate go run ../../tools/generate-version.go

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/ui"
	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
)

const (
	UserSessionIdCookie = "UserSessionId"
)

const (
	LocalsUserLoggedIn = "loggedIn"
	LocalsUser         = "user"
	LocalsUserToken    = "userToken"
)

const (
	RenderIsLoggedIn = "loggedIn"
	RenderUser       = "user"
)

const (
	RoleAdmin = "admin"
)

type handler struct {
	userStore users.UserStore
}

func main() {
	flags := flag.NewFlagSet("", flag.PanicOnError)

	resourcesDir := flags.String("resourcesDir", "", "path to resources dir")

	flags.Parse(os.Args[1:])

	var engine *html.Engine
	var staticFilesFs http.FileSystem

	if *resourcesDir != "" {
		if _, err := os.Stat(*resourcesDir); errors.Is(err, os.ErrNotExist) {
			panic(fmt.Sprintf("folder `%s` does not exist", *resourcesDir))
		}

		engine = html.New(filepath.Join(*resourcesDir, "views"), ".html")
		engine.Reload(true)

		staticFilesFs = http.Dir(filepath.Join(*resourcesDir, "static"))
	} else {
		engine = ui.CreateEmbeddedEngine()
		staticFilesFs = ui.GetStaticFilesFs()
	}

	handler := handler{
		userStore: users.BuildInMemoryStore(),
	}

	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: errorHandler,
	})

	app.Use(handler.setLoggedIn)

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: staticFilesFs,
	}))

	app.Get("/", handler.checkLoginOrRedirectToLoginPage, handler.renderIndex)
	app.Get("/admin", handler.checkLoginOrRedirectToLoginPage, handler.enforceRole(RoleAdmin), handler.renderIndex)

	app.Get("/login", handler.renderLogin)
	app.Post("/login", handler.doLogin)

	app.Get("/logout", handler.doLogout)
	app.Post("/logout", handler.doLogout)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("errors/404", handler.withDefault(c))
	})

	log.Fatal(app.Listen(":8080"))
}

func (h *handler) withDefault(c *fiber.Ctx, binds ...fiber.Map) fiber.Map {
	target := fiber.Map{}

	loggedIn := c.Locals(LocalsUserLoggedIn) == true
	target[RenderIsLoggedIn] = loggedIn

	if loggedIn {
		user := c.Locals(LocalsUser)

		if _, ok := user.(*users.ViewableUser); ok {
			target[RenderUser] = user
		}
	}

	for _, b := range binds {
		for k, v := range b {
			target[k] = v
		}
	}

	return target
}

func (h *handler) renderAdmin(c *fiber.Ctx) error {
	// currently only a dummy implementation for testing role enforcement
	return c.Render("index", h.withDefault(c), "layouts/main")
}

func (h *handler) renderIndex(c *fiber.Ctx) error {
	return c.Render("index", h.withDefault(c), "layouts/main")
}

func (h *handler) doLogout(c *fiber.Ctx) error {
	c.ClearCookie(UserSessionIdCookie)

	if token, ok := c.Locals(LocalsUserToken).(string); ok {
		h.userStore.InvalidateToken(token)
	}

	return c.Redirect("/login")
}

func (h *handler) renderLogin(c *fiber.Ctx) error {
	return c.Render("login", h.withDefault(c), "layouts/main")
}

func (h *handler) doLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	cookie, err := h.userStore.TryLogin(username, password)

	if errors.Is(err, users.InvalidUsernameOrPassword) {
		return c.Render("login", h.withDefault(c, fiber.Map{
			"invalidLogin": true,
		}), "layouts/main")
	}

	c.Cookie(&fiber.Cookie{
		Name:  UserSessionIdCookie,
		Value: cookie,
	})

	return c.Redirect("/")
}

func (h *handler) checkLoginOrRedirectToLoginPage(c *fiber.Ctx) error {
	loggedIn := c.Locals(LocalsUserLoggedIn)

	if loggedIn == true {
		return c.Next()
	}

	return c.Redirect("/login")
}

func (h *handler) enforceRole(role string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if user, ok := c.Locals(LocalsUser).(*users.ViewableUser); ok {
			if !hasRole(user, RoleAdmin) && !hasRole(user, role) {
				return c.Status(fiber.StatusForbidden).Render("errors/401", h.withDefault(c))
			}
		}

		return c.Next()
	}
}

func hasRole(user *users.ViewableUser, role string) bool {
	for _, r := range user.Roles {
		if r == role {
			return true
		}
	}

	return false
}

func (h *handler) setLoggedIn(c *fiber.Ctx) error {
	sessionId := c.Cookies(UserSessionIdCookie)

	if valid, user := h.userStore.ValidateToken(sessionId); valid {
		c.Locals(LocalsUserLoggedIn, true)
		c.Locals(LocalsUser, user)
		c.Locals(LocalsUserToken, sessionId)
	} else {
		c.Locals(LocalsUserLoggedIn, false)
	}

	return c.Next()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if code == fiber.StatusInternalServerError {
		log.Println(err)
	}

	err = c.Status(code).Render(fmt.Sprintf("errors/%d", code), fiber.Map{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return nil
}
