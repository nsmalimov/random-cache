package worker

import (
	"fmt"
	"random-cache/internal/cache"
)

type LastTwoElemsFromCache struct {
	FirstElem  interface{} `json:"first_elem"`
	SecondElem interface{} `json:"second_elem"`
}

type Worker struct {
	cache *cache.Cache
}

func New() *Worker {
	return &Worker{
		cache: cache.New(),
	}
}

func (w *Worker) LastTwoElemsFromCache() (lastTwoElemsFromCache LastTwoElemsFromCache, err error) {
	cacheSize := w.cache.Size()

	if cacheSize < 2 {
		err = fmt.Errorf("No 2 elem in cache, current size: %d", cacheSize)
		return
	}

	if firstElem, ok := w.cache.ItemByIndex(cacheSize - 1); !ok {
		err = fmt.Errorf("No elem with index: %d in cache", cacheSize-1)
		return
	} else {
		lastTwoElemsFromCache.FirstElem = firstElem
	}

	if secondElem, ok := w.cache.ItemByIndex(cacheSize - 2); !ok {
		err = fmt.Errorf("No elem with index: %d in cache", cacheSize-2)
	} else {
		lastTwoElemsFromCache.SecondElem = secondElem
	}

	return
}

func (w *Worker) Close() {

}

// start gen and add new string to cache
func (w *Worker) Init() {

}
