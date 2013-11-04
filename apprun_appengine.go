// +build appengine

package beego

import (
	"net/http"
)

func (app *App) Run() {
	http.Handle("/", app.Handlers)
}
