package core

import (
	"sort"
)

const EvictionPoolMaxSize = 16

var evictionPool = newEvictionPool(0)

type PoolItem struct {
	key              string
	lastAccessedTime uint32
}

type EvictionPool struct {
	items     []*PoolItem
	keyMapper map[string]*PoolItem
}

func newEvictionPool(size int) *EvictionPool {
	return &EvictionPool{
		items:     make([]*PoolItem, 0, size),
		keyMapper: make(map[string]*PoolItem),
	}

}

type IdleTime []*PoolItem

func (a IdleTime) Len() int {
	return len(a)
}

func (a IdleTime) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a IdleTime) Less(i, j int) bool {
	return getIdleTime(a[i].lastAccessedTime) > getIdleTime(a[j].lastAccessedTime)
}

// TODO : Improve it better
func (ep *EvictionPool) push(key string, lastAccessedTime uint32) {
	if _, ok := ep.keyMapper[key]; ok {
		return
	}
	item := &PoolItem{
		key:              key,
		lastAccessedTime: lastAccessedTime,
	}
	if len(ep.items) < EvictionPoolMaxSize {

		ep.items = append(ep.items, item)
		// Slow method
		sort.Sort(IdleTime(ep.items))

	} else if getIdleTime(lastAccessedTime) > getIdleTime(ep.items[len(ep.items)-1].lastAccessedTime) {
		ep.items = ep.items[1:]
		ep.keyMapper[key] = item
		deletedItem := ep.items[0]
		delete(ep.keyMapper, deletedItem.key)
		ep.items = append(ep.items, item)

	}

}

func (ep *EvictionPool) pop() *PoolItem {

	if len(ep.items) == 0 {
		return nil
	}
	ep.items = ep.items[1:]
	deletedItem := ep.items[0]
	delete(ep.keyMapper, deletedItem.key)
	return deletedItem
}
