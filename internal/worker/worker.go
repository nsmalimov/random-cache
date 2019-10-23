package worker

import (
	"fmt"
	"time"

	"random-cache/internal/cache"
	"random-cache/internal/config"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LastTwoElemsFromCache struct {
	Elems []interface{} `json:"elems"`
}

type Worker struct {
	cache    *cache.Cache
	cfg      *config.Config
	log      *logrus.Logger
	stopChan chan bool
}

func New(config *config.Config, logger *logrus.Logger) *Worker {
	return &Worker{
		cache:    cache.New(),
		log:      logger,
		cfg:      config,
		stopChan: make(chan bool),
	}
}

func (w *Worker) ElemsFromCache() (lastTwoElemsFromCache LastTwoElemsFromCache, err error) {
	cacheSize := w.cache.Size()

	if cacheSize < w.cfg.HowMuchLastElemsFromCacheNeedReturn {
		err = fmt.Errorf("no %d elem in cache, current size: %d", w.cfg.HowMuchLastElemsFromCacheNeedReturn,
			cacheSize)
		return
	}

	for i := 0; i < w.cfg.HowMuchLastElemsFromCacheNeedReturn; i++ {
		elem, ok := w.cache.ItemByIndex(i)

		// по сути невозможный кейс. но проверку оставляю
		if !ok {
			err = fmt.Errorf("no elem with index: %d in cache", i)
			return
		}

		lastTwoElemsFromCache.Elems = append(lastTwoElemsFromCache.Elems, elem)
	}

	return
}

func (w *Worker) Close() {
	w.stopChan <- true
}

func (w *Worker) getString() string {
	return uuid.New().String()[:w.cfg.LenStringForAddToCache]
}

// start gen and add new string to cache
func (w *Worker) Init() {
	go func() {
		for {
			select {
			case <-w.stopChan:
				return
			case <-time.After(time.Second * time.Duration(w.cfg.FrequencyAddToCacheSec)):
				s := w.getString()
				w.cache.AddItem(s)
				w.log.Infof("Elem: %+v was added to cache, current size: %d", s, w.cache.Size())
			}
		}
	}()
}
