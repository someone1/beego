// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Server   *http.Server
}

// NewApp returns a new beego application.
func NewApp() *App {
	cr := NewControllerRegister()
	app := &App{Handlers: cr, Server: &http.Server{}}
	return app
}

// Run beego application.
func (app *App) Run() {
	http.Handle("/", app.Handlers)
}
