// +build !appengine

package beego

import (
	"os"
	"path"
	"runtime"
)

func init_platform() {
	os.Chdir(path.Dir(os.Args[0])) //Is this required?
	AppPath = path.Dir(os.Args[0])

	runtime.GOMAXPROCS(runtime.NumCPU())
}
