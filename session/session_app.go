// +build appengine

package session

import (
	"net/http"
	"net/url"

	"appengine"
)

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
