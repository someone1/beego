beegae session
==============

This is based off of the original [session module](https://github.com/astaxie/beego/tree/master/session) as part of the beego project and can be used as part of the beegae project or as a standalone session manager. [read more here](http://beego.me/docs/mvc/controller/session.md)

This includes (as of now) only a single session store capable of working on AppEngine (`SessionProvider = "appengine"`)

A few gotchas:

1. There is no automatic garbage collection! You will have to create a cron job and a custom handler to periodically call on the garbage collection functions.
2. `SessionAll` will always return 0. `Count` queries are limited to 1000 entities and so we cannot reliably get a count. As such, this function was not implemented.
3. A few methods deviate from the original beego API specification. Specifically, an `appengine.Context` object is a new parameter for `SessionExist`, `SessionRead`, `SessionRegenerate`, `SessionDestroy`, and `SessionGC`. If you are using beegae or use the session manager provided, you do not have to worry about these details.
4. `GetProvider` was not implemented (this should have little to no impact)

Example Garbage Collection using **beegae**:

First, create a new controller:

```go
package controllers

import "github.com/astaxie/beegae"

type GCController struct {
	beegae.Controller
}

func (this *GCController) Get() {
	beegae.GlobalSessions.GC(this.AppEngineCtx)
}
```

Second, register your controller to a URL Path in your applications `init` function:

```go
func init() {
	// Register other routers/handlers here
	// ...

	// Register new handler for sessiong garbage collection
	beegae.Router("/_session_gc", &controllers.GCController)

	beegae.Run()
}
```

<<<<<<< HEAD
Finally, add an entry to your cron.yaml file:
=======
		func init() {
			globalSessions, _ = session.NewManager("file",`{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./tmp"}`)
			go globalSessions.GC()
		}
>>>>>>> 03080b3ef280c2e8ada5436d670dcc0772dd2a5b

```yaml
cron:
- description: daily session garbage collection
  url: /_session_gc
  schedule: every day 00:00
```

<<<<<<< HEAD
You can also add security to this (and any) URL by requiring an Admin login for the URL in your app.yaml:

```yaml
handlers:
- url: /_session_gc
  login: admin
  script: _go_app

- url: /.*
  script: _go_app
```
=======
		func init() {
			globalSessions, _ = session.NewManager("redis", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:6379,100,astaxie"}`)
			go globalSessions.GC()
		}
		
* Use **MySQL** as provider, the last param is the DSN, learn more from [mysql](https://github.com/go-sql-driver/mysql#dsn-data-source-name):

		func init() {
			globalSessions, _ = session.NewManager(
				"mysql", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"username:password@protocol(address)/dbname?param=value"}`)
			go globalSessions.GC()
		}

* Use **Cookie** as provider:

		func init() {
			globalSessions, _ = session.NewManager(
				"cookie", `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`)
			go globalSessions.GC()
		}


Finally in the handlerfunc you can use it like this

	func login(w http.ResponseWriter, r *http.Request) {
		sess := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		username := sess.Get("username")
		fmt.Println(username)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
		} else {
			fmt.Println("username:", r.Form["username"])
			sess.Set("username", r.Form["username"])
			fmt.Println("password:", r.Form["password"])
		}
	}


## How to write own provider?

When you develop a web app, maybe you want to write own provider because you must meet the requirements.

Writing a provider is easy. You only need to define two struct types 
(Session and Provider), which satisfy the interface definition. 
Maybe you will find the **memory** provider is a good example.

	type SessionStore interface {
		Set(key, value interface{}) error     //set session value
		Get(key interface{}) interface{}      //get session value
		Delete(key interface{}) error         //delete session value
		SessionID() string                    //back current sessionID
		SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
		Flush() error                         //delete all data
	}
	
	type Provider interface {
		SessionInit(gclifetime int64, config string) error
		SessionRead(sid string) (SessionStore, error)
		SessionExist(sid string) bool
		SessionRegenerate(oldsid, sid string) (SessionStore, error)
		SessionDestroy(sid string) error
		SessionAll() int //get all active session
		SessionGC()
	}


## LICENSE

BSD License http://creativecommons.org/licenses/BSD/
>>>>>>> 03080b3ef280c2e8ada5436d670dcc0772dd2a5b
