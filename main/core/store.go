package core

import "time"

// Starting the project just with hashmap. Will optimize later
var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObject(value interface{}, durationMs int64, objType uint8, objEncoding uint8) *Obj {
	var expiry int64 = -1
	if durationMs > 0 {
		expiry = durationMs + time.Now().UnixMilli()
	}
	return &Obj{
		TypeEncoding: objType | objEncoding,
		Value:        value,
		Expiry:       expiry,
	}

}

func Put(key string, value *Obj) {
	store[key] = value
}

func Get(key string) *Obj {
	return store[key]
}

func Delete(key string) bool {
	if _, ok := store[key]; ok {
		delete(store, key)
		return true
	}
	return false
}
