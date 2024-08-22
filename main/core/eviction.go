package core

import "myredis/config"

func evictKeys() {
	switch config.EvictionStrategy {
	case "simple_first":
		evictFirst()
	}

}

func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}
