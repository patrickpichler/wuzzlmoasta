package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/template/html"
)

//go:embed app
var embeddedFiles embed.FS

func CreateEmbeddedEngine() *html.Engine {
  targetFs := mustFs(fs.Sub(embeddedFiles, "app/views"))

	return html.NewFileSystem(http.FS(targetFs), ".html")
}

func GetStaticFilesFs() http.FileSystem {
  targetFs := mustFs(fs.Sub(embeddedFiles, "app/static"))

  return http.FS(targetFs)
}

func mustFs(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}

	return fs
}
