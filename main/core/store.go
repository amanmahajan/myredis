package core

import (
	"myredis/config"
	"time"
)

// Starting the project just with hashmap. Will optimize later
var store map[string]*Obj
var expires map[*Obj]uint64

func init() {
	store = make(map[string]*Obj)
	expires = make(map[*Obj]uint64)
}

func setExpiry(obj *Obj, expireTime int64) {
	expires[obj] = uint64(expireTime) + uint64(time.Now().UnixMilli())
}

func NewObject(value interface{}, durationMs int64, objType uint8, objEncoding uint8) *Obj {

	obj := &Obj{
		TypeEncoding:   objType | objEncoding,
		Value:          value,
		LastAccessedAt: getCurrentClock(),
	}

	if durationMs > 0 {
		setExpiry(obj, durationMs)
	}

	return obj

}

func Put(key string, value *Obj) {

	if len(store) >= config.KeyLimits {
		evictKeys()
	}
	value.LastAccessedAt = getCurrentClock()
	if KeySpaceStats[0] == nil {
		KeySpaceStats[0] = make(map[string]int, 0)
	}
	store[key] = value
	KeySpaceStats[0]["Keys"]++
}

func Get(key string) *Obj {
	val := store[key]
	if val != nil {
		if hasExpired(val) {
			Delete(key)
			return nil

		}
		val.LastAccessedAt = getCurrentClock()
	}
	return val

}

func Delete(key string) bool {
	if obj, ok := store[key]; ok {
		delete(store, key)
		delete(expires, obj)
		KeySpaceStats[0]["Keys"]--
		return true
	}

	return false
}
