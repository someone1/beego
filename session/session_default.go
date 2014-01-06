// +build !appengine

package session

import (
	"net/http"
	"net/url"
	"time"
)

type Provider interface {
	SessionInit(gclifetime int64, config string) error
	SessionRead(sid string) (SessionStore, error)
	SessionExist(sid string) bool
	SessionRegenerate(oldsid, sid string) (SessionStore, error)
	SessionDestroy(sid string) error
	SessionAll() int //get all active session
	SessionGC()
}

//get Session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session SessionStore) {
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId(r)
		session, _ = manager.provider.SessionRead(sid)
		cookie = &http.Cookie{Name: manager.config.CookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			Secure:   manager.config.Secure}
		if manager.config.Maxage >= 0 {
			cookie.MaxAge = manager.config.Maxage
		}
		if manager.config.EnableSetCookie {
			http.SetCookie(w, cookie)
		}
		r.AddCookie(cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		if manager.provider.SessionExist(sid) {
			session, _ = manager.provider.SessionRead(sid)
		} else {
			sid = manager.sessionId(r)
			session, _ = manager.provider.SessionRead(sid)
			cookie = &http.Cookie{Name: manager.config.CookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: true,
				Secure:   manager.config.Secure}
			if manager.config.Maxage >= 0 {
				cookie.MaxAge = manager.config.Maxage
			}
			if manager.config.EnableSetCookie {
				http.SetCookie(w, cookie)
			}
			r.AddCookie(cookie)
		}
	}
	return
}

//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.config.CookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.provider.SessionGC()
	time.AfterFunc(time.Duration(manager.config.Gclifetime)*time.Second, func() { manager.GC() })
}

func (manager *Manager) SessionRegenerateId(w http.ResponseWriter, r *http.Request) (session SessionStore) {
	sid := manager.sessionId(r)
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil && cookie.Value == "" {
		//delete old cookie
		session, _ = manager.provider.SessionRead(sid)
		cookie = &http.Cookie{Name: manager.config.CookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			Secure:   manager.config.Secure,
		}
	} else {
		oldsid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRegenerate(oldsid, sid)
		cookie.Value = url.QueryEscape(sid)
		cookie.HttpOnly = true
		cookie.Path = "/"
	}
	if manager.config.Maxage >= 0 {
		cookie.MaxAge = manager.config.Maxage
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	return
}
