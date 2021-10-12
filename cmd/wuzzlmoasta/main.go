package main

//go:generate go run ../../tools/generate-version.go

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/ui"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "")
}

func main() {
	flags := flag.NewFlagSet("", flag.PanicOnError)

	resourcesDir := flags.String("resourcesDir", "", "path to resources dir")

	flags.Parse(os.Args[1:])

	e := echo.New()

	var webUI ui.UI

	if *resourcesDir != "" {
		if _, err := os.Stat(*resourcesDir); errors.Is(err, os.ErrNotExist) {
			panic(fmt.Sprintf("folder `%s` does not exist", *resourcesDir))
		}

		webUI = ui.NewLive(*resourcesDir)
	} else {
		webUI = ui.New()
	}

	e.Renderer = webUI

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		errorPage := fmt.Sprintf("%d.html", code)

		if err := c.Render(code, errorPage, nil); err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Error(err)
	}

	assetHandler := http.FileServer(http.FS(webUI.StaticAssets()))

	e.GET("/", index)
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	e.Logger.Fatal(e.Start(":8080"))
}
