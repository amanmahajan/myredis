package core

import (
	"log"
	"time"
)

const sampleKeysLimit = 20

// Redis Sampling method to delete he expired keys.
func expireSampleKeys() float32 {

	count := 0
	deletedKeys := 0

	for key, val := range store {
		if count > sampleKeysLimit {
			break
		}
		if val.Expiry <= time.Now().UnixMilli() {

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
