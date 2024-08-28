package core

import (
	"myredis/config"
	"time"
)

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

	if len(store) > config.KeyLimits {
		evictAllKeyRandom()
	}
	if KeySpaceStats[0] == nil {
		KeySpaceStats[0] = make(map[string]int, 0)
	}
	store[key] = value
	KeySpaceStats[0]["Keys"]++
}

func Get(key string) *Obj {
	val := store[key]
	if val != nil {
		if val.Expiry != -1 && val.Expiry >= time.Now().UnixMilli() {
			Delete(key)
			return nil

		}
	}
	return val

}

func Delete(key string) bool {
	if _, ok := store[key]; ok {
		delete(store, key)
		KeySpaceStats[0]["Keys"]--
		return true
	}

	return false
}
