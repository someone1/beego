// +build appengine

package session

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"

	"time"
)

var appenginepvdr = &AppEngineProvider{}

type AppEngineSessionStore struct {
	c           appengine.Context
	sid         string
	lock        sync.RWMutex
	dirty       bool
	maxlifetime int64
	bss_entity  *BeegoSessionStore
	values      map[interface{}]interface{}
}

type BeegoSessionStore struct {
	SessionData  []byte
	SessionStart time.Time
}

func (st *AppEngineSessionStore) Set(key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values[key] = value
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) Get(key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if v, ok := st.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *AppEngineSessionStore) Delete(key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	delete(st.values, key)
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) Flush() error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values = make(map[interface{}]interface{})
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) SessionID() string {
	return st.sid
}

func (st *AppEngineSessionStore) updatestore() {
	b, err := encodeGob(st.values)
	if err != nil {
		st.c.Errorf("error encoding session data: %v", err)
		return
	}

	done := make(chan bool, 2)

	if st.bss_entity == nil {
		st.bss_entity = &BeegoSessionStore{SessionStart: time.Now()}
	}

	st.bss_entity.SessionData = b

	go func() {
		k := datastore.NewKey(c, "BeegoSessionStore", st.sid, 0, nil)
		if ds_err := datastore.Put(st.c, k, st.bss_entity); ds_err != nil {
			st.c.Errorf("error saving session data to datastore: %v", ds_err)
		}
		done <- true
	}()

	go func() {
		mem_err := memcache.Set(st.c, &memcache.Item{
			Key:        st.sid,
			Value:      st.bss_entity.SessionData,
			Expiration: (time.Duration(st.maxlifetime) * time.Second) - time.Since(st.bss_entity.SessionStart),
		})
		if mem_err != nil {
			st.c.Errorf("error saving session data to memcache: %v", mem_err)
		}
		done <- true
	}()

	mem_wait, ds_wait := <-done, <-done
}

func (st *AppEngineSessionStore) SessionRelease() {
	//Always expected to be called to save session data
	if st.dirty {
		st.updatestore()
	}
}

type AppEngineProvider struct {
	maxlifetime int64
	savePath    string
}

func (mp *AppEngineProvider) SessionInit(maxlifetime int64, savePath string) error {
	mp.maxlifetime = maxlifetime
	mp.savePath = savePath
	return nil
}

func (mp *AppEngineProvider) SessionRead(sid string) (SessionStore, error) {
	panic("Who called me? I'm not used for the AppEngine Session backend!")
}

func (mp *AppEngineProvider) getsession(sid string, c appengine.Context) *BeegoSessionStore {
	in_cache := false
	e := new(BeegoSessionStore)
	if item, err := memcache.Get(c, sid); err == memcache.ErrCacheMiss {
		//This is ok!
	} else if err != nil {
		c.Errorf("error getting session data from memcache: %v", err)
	} else {
		in_cache = true
		e.SessionData = item.Value
		e.SessionStart = time.Now().Add(-(time.Duration(mp.maxlifetime) * time.Second) - item.Expiration)
	}

	if !in_cache {
		k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
		if ds_err := datastore.Get(c, k, e); ds_err == datastore.ErrNoSuchEntity {
			e.SessionStart = time.Now()
		} else if ds_err != nil {
			c.Errorf("error getting session data from datastore: %v", err)
			e.SessionStart = time.Now()
		}
	}
	return e
}

func (mp *AppEngineProvider) SessionExist(sid string, c appengine.Context) bool {
	k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
	if ds_err := datastore.Get(c, k, e); ds_err == datastore.ErrNoSuchEntity {
		return false
	} else if ds_err != nil {
		c.Errorf("error while checking existence of session data from datastore: %v", ds_err)
		return false
	} else {
		return true
	}
}

func (mp *AppEngineProvider) SessionRead(sid string, c appengine.Context) (SessionStore, error) {
	e := getsession(sid, c)
	var kv map[interface{}]interface{}

	if len(e.SessionData) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(e.SessionData)
		if err != nil {
			return nil, err
		}
	}
	rs := &AppEngineSessionStore{c: c, sid: sid, values: kv, maxlifetime: maxlifetime, dirty: false, bss_entity: e}
	return rs, nil
}

func (mp *AppEngineProvider) SessionRegenerate(oldsid, sid string) (SessionStore, error) {
	panic("Who called me? I'm not used for the AppEngine Session backend!")
}

func (mp *AppEngineProvider) SessionRegenerate(oldsid, sid string, c appengine.Context) (SessionStore, error) {
	e := getsession(sid, c)
	var kv map[interface{}]interface{}

	if len(e.SessionData) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(e.SessionData)
		if err != nil {
			return nil, err
		}
	}
	rs := &AppEngineSessionStore{c: c, sid: sid, values: kv, maxlifetime: maxlifetime, dirty: false, bss_entity: e}
	return rs, nil
}

func (mp *AppEngineProvider) SessionDestroy(sid string, c appengine.Context) error {
	go func() {
		k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
		if ds_err := datastore.Delete(c, k); ds_err != nil {
			c.Errorf("error deleting session data from datastore: %v", ds_err)
		}
		done <- true
	}()

	go func() {
		mem_err := memcache.Delete(c, sid)
		if mem_err != nil {
			c.Errorf("error deleting session data from memcache: %v", mem_err)
		}
		done <- true
	}()

	return nil
}

func (mp *AppEngineProvider) SessionGC(c appengine.Context) {
	q := datastore.NewQuery("BeegoSessionStore").Filter("SessionStart <", time.Now().Unix()-mp.maxlifetime).KeysOnly()

	keys, err := q.GetAll(c, nil)
	if err != nil {
		c.Errorf("error querying session data from datastore: %v", err)
	}

	for key := range keys {
		mp.SessionDestroy(key.StringID(), c)
	}
	return
}

func (mp *AppEngineProvider) SessionAll(c appengine.Context) int {
	total, err := datastore.NewQuery("BeegoSessionStore").KeysOnly().Count(c)
	if err != nil {
		return 0
	}
	return total
}

func init() {
	Register("appengine", appenginepvdr)
}
