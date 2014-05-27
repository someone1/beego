// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package beegae

import (
	"github.com/astaxie/beegae/context"
	"net/http"
)

// FilterFunc defines filter function type.
type FilterFunc func(*context.Context)

// App defines beego application with a new PatternServeMux.
type App struct {
	Handlers *ControllerRegistor
}

// NewApp returns a new beego application.
func NewApp() *App {
	cr := NewControllerRegistor()
	app := &App{Handlers: cr}
	return app
}

// Router adds a url-patterned controller handler.
// The path argument supports regex rules and specific placeholders.
// The c argument needs a controller handler implemented beego.ControllerInterface.
// The mapping methods argument only need one string to define custom router rules.
// usage:
//  simple router
//  beego.Router("/admin", &admin.UserController{})
//  beego.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  beego.Router(“/api/:id([0-9]+)“, &controllers.RController{})
//
//  custom rules
//  beego.Router("/api/list",&RestController{},"*:ListFood")
//  beego.Router("/api/create",&RestController{},"post:CreateFood")
//  beego.Router("/api/update",&RestController{},"put:UpdateFood")
//  beego.Router("/api/delete",&RestController{},"delete:DeleteFood")
func (app *App) Router(path string, c ControllerInterface, mappingMethods ...string) *App {
	app.Handlers.Add(path, c, mappingMethods...)
	return app
}

// AutoRouter adds beego-defined controller handler.
// if beego.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func (app *App) AutoRouter(c ControllerInterface) *App {
	app.Handlers.AddAuto(c)
	return app
}

// AutoRouterWithPrefix adds beego-defined controller handler with prefix.
// if beego.AutoPrefix("/admin",&MainContorlller{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func (app *App) AutoRouterWithPrefix(prefix string, c ControllerInterface) *App {
	app.Handlers.AddAutoPrefix(prefix, c)
	return app
}

// add router for Get method
func (app *App) Get(rootpath string, f FilterFunc) *App {
	app.Handlers.Get(rootpath, f)
	return app
}

// add router for Post method
func (app *App) Post(rootpath string, f FilterFunc) *App {
	app.Handlers.Post(rootpath, f)
	return app
}

// add router for Put method
func (app *App) Put(rootpath string, f FilterFunc) *App {
	app.Handlers.Put(rootpath, f)
	return app
}

// add router for Delete method
func (app *App) Delete(rootpath string, f FilterFunc) *App {
	app.Handlers.Delete(rootpath, f)
	return app
}

// add router for Options method
func (app *App) Options(rootpath string, f FilterFunc) *App {
	app.Handlers.Options(rootpath, f)
	return app
}

// add router for Head method
func (app *App) Head(rootpath string, f FilterFunc) *App {
	app.Handlers.Head(rootpath, f)
	return app
}

// add router for Patch method
func (app *App) Patch(rootpath string, f FilterFunc) *App {
	app.Handlers.Patch(rootpath, f)
	return app
}

// add router for Patch method
func (app *App) Any(rootpath string, f FilterFunc) *App {
	app.Handlers.Any(rootpath, f)
	return app
}

// add router for http.Handler
func (app *App) Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	app.Handlers.Handler(rootpath, h, options...)
	return app
}

// UrlFor creates a url with another registered controller handler with params.
// The endpoint is formed as path.controller.name to defined the controller method which will run.
// The values need key-pair data to assign into controller method.
func (app *App) UrlFor(endpoint string, values ...string) string {
	return app.Handlers.UrlFor(endpoint, values...)
}

// [Deprecated] use InsertFilter.
// Filter adds a FilterFunc under pattern condition and named action.
// The actions contains BeforeRouter,AfterStatic,BeforeExec,AfterExec and FinishRouter.
func (app *App) Filter(pattern, action string, filter FilterFunc) *App {
	app.Handlers.AddFilter(pattern, action, filter)
	return app
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// beego.BeforeRouter, beego.AfterStatic, beego.BeforeExec, beego.AfterExec and beego.FinishRouter.
func (app *App) InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	app.Handlers.InsertFilter(pattern, pos, filter)
	return app
}

// SetViewsPath sets view directory path in beego application.
// it returns beego application self.
func (app *App) SetViewsPath(path string) *App {
	ViewsPath = path
	return app
}

// SetStaticPath sets static directory path and proper url pattern in beego application.
// if beego.SetStaticPath("static","public"), visit /static/* to load static file in folder "public".
// it returns beego application self.
func (app *App) SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return app
}

// DelStaticPath removes the static folder setting in this url pattern in beego application.
// it returns beego application self.
func (app *App) DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return app
}

// Run beego application.
func (app *App) Run() {
	http.Handle("/", app.Handlers)
}
