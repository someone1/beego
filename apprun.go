// +build !appengine

package beego

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"time"
)

func (app *App) Run() {
	addr := HttpAddr

	if HttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	}

	BeeLogger.Info("Running on %s", addr)

	var (
		err error
		l   net.Listener
	)

	if UseFcgi {
		if HttpPort == 0 {
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			BeeLogger.Critical("Listen: ", err)
		}
		err = fcgi.Serve(l, app.Handlers)
	} else {
		if EnableHotUpdate {
			server := &http.Server{
				Handler:      app.Handlers,
				ReadTimeout:  time.Duration(HttpServerTimeOut) * time.Second,
				WriteTimeout: time.Duration(HttpServerTimeOut) * time.Second,
			}
			laddr, err := net.ResolveTCPAddr("tcp", addr)
			if nil != err {
				BeeLogger.Critical("ResolveTCPAddr:", err)
			}
			l, err = GetInitListner(laddr)
			theStoppable = newStoppable(l)
			err = server.Serve(theStoppable)
			theStoppable.wg.Wait()
			CloseSelf()
		} else {
			s := &http.Server{
				Addr:         addr,
				Handler:      app.Handlers,
				ReadTimeout:  time.Duration(HttpServerTimeOut) * time.Second,
				WriteTimeout: time.Duration(HttpServerTimeOut) * time.Second,
			}
			if HttpTLS {
				err = s.ListenAndServeTLS(HttpCertFile, HttpKeyFile)
			} else {
				err = s.ListenAndServe()
			}
		}
	}

	if err != nil {
		BeeLogger.Critical("ListenAndServe: ", err)
		time.Sleep(100 * time.Microsecond)
	}
}
