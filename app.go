package beego

import (
	"github.com/astaxie/beego/context"
)

type FilterFunc func(*context.Context)

type App struct {
	Handlers *ControllerRegistor
}

// New returns a new PatternServeMux.
func NewApp() *App {
	cr := NewControllerRegistor()
	app := &App{Handlers: cr}
	return app
}

func (app *App) Router(path string, c ControllerInterface, mappingMethods ...string) *App {
	app.Handlers.Add(path, c, mappingMethods...)
	return app
}

func (app *App) AutoRouter(c ControllerInterface) *App {
	app.Handlers.AddAuto(c)
	return app
}

func (app *App) UrlFor(endpoint string, values ...string) string {
	return app.Handlers.UrlFor(endpoint, values...)
}
func (app *App) Filter(pattern, action string, filter FilterFunc) *App {
	app.Handlers.AddFilter(pattern, action, filter)
	return app
}

func (app *App) InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	app.Handlers.InsertFilter(pattern, pos, filter)
	return app
}

func (app *App) SetViewsPath(path string) *App {
	ViewsPath = path
	return app
}

func (app *App) SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return app
}

func (app *App) DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return app
}
