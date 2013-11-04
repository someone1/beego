// +build !appengine

package beego

import (
	"html/template"
	"os"
	"path"
	"runtime"
)

func init() {
	os.Chdir(path.Dir(os.Args[0]))
	BeeApp = NewApp()
	AppPath = path.Dir(os.Args[0])
	StaticDir = make(map[string]string)
	TemplateCache = make(map[string]*template.Template)
	HttpAddr = ""
	HttpPort = 8080
	AppName = "beego"
	RunMode = "dev" //default runmod
	AutoRender = true
	RecoverPanic = true
	PprofOn = false
	ViewsPath = path.Join(AppPath, "views")
	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
	SessionHashFunc = "sha1"
	SessionHashKey = "beegoserversessionkey"
	SessionCookieLifeTime = 3600
	UseFcgi = false
	MaxMemory = 1 << 26 //64MB
	EnableGzip = false
	StaticDir["/static"] = "static"
	AppConfigPath = path.Join(AppPath, "conf", "app.conf")
	HttpServerTimeOut = 0
	ErrorsShow = true
	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0
	TemplateLeft = "{{"
	TemplateRight = "}}"
	BeegoServerName = "beegoServer"
	ParseConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())
}
