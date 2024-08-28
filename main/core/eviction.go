package core

import "myredis/config"

func evictKeys() {
	switch config.EvictionStrategy {
	case "simple_first":
		evictFirst()
	case "allkeys-random":
		evictAllKeyRandom()
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
