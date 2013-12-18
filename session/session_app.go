// +build appengine

package session

import (
	"net/http"
	"net/url"
  "time"

	"appengine"
)

type Provider interface {
	SessionInit(maxlifetime int64, savePath string) error
	SessionRead(sid string, c appengine.Context) (SessionStore, error)
	SessionExist(sid string, c appengine.Context) bool
	SessionRegenerate(oldsid, sid string, c appengine.Context) (SessionStore, error)
	SessionDestroy(sid string, c appengine.Context) error
	SessionAll() int //get all active session
	SessionGC(c appengine.Context)
}

//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
  var c = appengine.NewContext(r)
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.provider.SessionDestroy(cookie.Value, c)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) SessionRegenerateId(w http.ResponseWriter, r *http.Request) (session SessionStore) {
  var c = appengine.NewContext(r)
	sid := manager.sessionId(r)
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil && cookie.Value == "" {
		//delete old cookie
		session, _ = manager.provider.SessionRead(sid, c)
		cookie = &http.Cookie{Name: manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			Secure:   manager.secure,
		}
	} else {
		oldsid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRegenerate(oldsid, sid, c)
		cookie.Value = url.QueryEscape(sid)
		cookie.HttpOnly = true
		cookie.Path = "/"
	}
	if manager.maxage >= 0 {
		cookie.MaxAge = manager.maxage
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	return
}


//get Session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session SessionStore) {
	var c = appengine.NewContext(r)
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId(r)
		session, _ = manager.provider.SessionRead(sid, c)
		cookie = &http.Cookie{Name: manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			Secure:   manager.secure}
		if manager.maxage >= 0 {
			cookie.MaxAge = manager.maxage
		}
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		if manager.provider.SessionExist(sid, c) {
			session, _ = manager.provider.SessionRead(sid, c)
		} else {
			sid = manager.sessionId(r)
			session, _ = manager.provider.SessionRead(sid, c)
			cookie = &http.Cookie{Name: manager.cookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: true,
				Secure:   manager.secure}
			if manager.maxage >= 0 {
				cookie.MaxAge = manager.maxage
			}
			http.SetCookie(w, cookie)
			r.AddCookie(cookie)
		}
	}
	return
}

func (manager *Manager) GC() {
	return
}
