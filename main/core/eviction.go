package core

import (
	"myredis/config"
	"time"
)

func evictKeys() {
	switch config.EvictionStrategy {
	case "simple_first":
		evictFirst()
	case "allkeys-random":
		evictAllKeyRandom()
	case "allkeys-lru":
		evictAllkeysLRU()
	}

}

func evictFirst() {
	for k := range store {
		Delete(k)
		return
	}
}

func evictAllKeyRandom() {

	totalKeysToEvict := int64(config.EvictionRatio * float64(config.KeyLimits))
	// Golang Dictionary has random insertion
	for k := range store {
		Delete(k)
		totalKeysToEvict--
		if totalKeysToEvict <= 0 {
			return
		}
	}
}

/*
•	time.Now().Unix() returns the current Unix timestamp in seconds.
•	The uint32() function converts this timestamp to an unsigned 32-bit integer.
•	The & 0x00FFFFFF operation masks the upper 8 bits, keeping only the lower 24 bits.
*/
func getCurrentClock() uint32 {
	return uint32(time.Now().Unix()) & 0x00FFFFFF
}

func getIdleTime(lastAccessTime uint32) uint32 {
	c := getCurrentClock()
	if c >= lastAccessTime {
		return c - lastAccessTime

	}
	return (0x00FFFFFF - lastAccessTime) + c

}

func populateEvictionPool() {
	sampleKeysLimit := 5
	for key := range store {
		evictionPool.push(key, store[key].LastAccessedAt)

		sampleKeysLimit--
		if sampleKeysLimit <= 0 {
			return
		}
	}

}

func evictAllkeysLRU() {
	populateEvictionPool()
	evictCount := int16(config.EvictionRatio * float64(config.KeyLimits))
	for i := 0; i < int(evictCount) && len(evictionPool.items) > 0; i++ {
		item := evictionPool.pop()
		if item == nil {
			return
		}
		Delete(item.key)
	}
}
