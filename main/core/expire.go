package core

import (
	"log"
	"time"
)

const sampleKeysLimit = 20

func hasExpired(obj *Obj) bool {
	val, ok := expires[obj]
	if !ok {
		return false
	}
	return val <= uint64(time.Now().UnixMilli())
}

func getExpiry(obj *Obj) (uint64, bool) {
	exp, ok := expires[obj]
	return exp, ok
}

// Redis Sampling method to delete he expired keys.
func expireSampleKeys() float32 {

	count := 0
	deletedKeys := 0

	for key, val := range store {
		if count > sampleKeysLimit {
			break
		}
		if hasExpired(val) {
			Delete(key)
			deletedKeys += 1
		}
		count++
	}
	return float32(deletedKeys) / float32(sampleKeysLimit)

}

func DeleteExpiredKey() {
	for {
		val := expireSampleKeys()
		if val < 0.25 {
			break
		}
	}
	log.Println("Deleted Expired keys. Total keys left", len(store))
}
