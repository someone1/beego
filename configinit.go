// +build !appengine

package beego

import (
	"os"
	"path/filepath"
	"runtime"
)

func init_platform() {
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(AppPath) //Is this required?

	runtime.GOMAXPROCS(runtime.NumCPU())
}
