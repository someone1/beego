// +build !appengine

package beego

import (
	"html/template"
	"os"
	"path"
	"runtime"
)

func init() {
	// create beeapp
	BeeApp = NewApp()

	// initialize default configurations
	os.Chdir(path.Dir(os.Args[0]))
	AppPath = path.Dir(os.Args[0])

	StaticDir = make(map[string]string)
	StaticDir["/static"] = "static"

	StaticExtensionsToGzip = []string{".css", ".js"}

	TemplateCache = make(map[string]*template.Template)

	// set this to 0.0.0.0 to make this app available to externally
	HttpAddr = ""
	HttpPort = 8080

	AppName = "beego"

	RunMode = "dev" //default runmod

	AutoRender = true

	RecoverPanic = true

	ViewsPath = "views"

	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
	SessionHashFunc = "sha1"
	SessionHashKey = "beegoserversessionkey"
	SessionCookieLifeTime = 0 //set cookie default is the brower life

	UseFcgi = false

	MaxMemory = 1 << 26 //64MB

	EnableGzip = false

	AppConfigPath = path.Join(AppPath, "conf", "app.conf")

	HttpServerTimeOut = 0

	ErrorsShow = true

	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0

	TemplateLeft = "{{"
	TemplateRight = "}}"

	BeegoServerName = "beegoServer"

	EnableAdmin = true
	AdminHttpAddr = "127.0.0.1"
	AdminHttpPort = 8088

	runtime.GOMAXPROCS(runtime.NumCPU())

	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
}
