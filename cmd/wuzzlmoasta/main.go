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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
)

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

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = c.Status(code).Render(fmt.Sprintf("errors/%d", code), fiber.Map{})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
			}

			return nil
		},
	})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: staticFilesFs,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

  app.Use(func(c *fiber.Ctx) error {
    return c.Status(404).Render("errors/404", fiber.Map{})
  })

	log.Fatal(app.Listen(":8080"))

}
