package ui

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed app
var embeddedFiles embed.FS

type UI interface {
	echo.Renderer

	StaticAssets() fs.FS
}

type inMemoryUI struct {
	fs        fs.FS
	templates *template.Template
}

func (t inMemoryUI) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t inMemoryUI) StaticAssets() fs.FS {
	return mustFs(fs.Sub(t.fs, "static"))
}

type liveUI struct {
	fs fs.FS
}

func (t liveUI) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templates := template.Must(template.ParseFS(t.fs, "views/*"))

	return templates.ExecuteTemplate(w, name, data)
}

func (t liveUI) StaticAssets() fs.FS {
	return mustFs(fs.Sub(t.fs, "static"))
}

func getFileSystem() fs.FS {
	return mustFs(fs.Sub(embeddedFiles, "app"))
}

func New() UI {
	fs := getFileSystem()

	return &inMemoryUI{
		fs:        fs,
		templates: template.Must(template.ParseFS(fs, "views/*")),
	}
}

func NewLive(path string) UI {
	fs := os.DirFS(path)

	return &liveUI{
		fs: fs,
	}
}

func mustFs(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}

	return fs
}
