package core

import (
	"errors"
	"push-api-cron/core/models/device"
	"sync"
	"time"
)

type Item struct {
	Value      []device.Device
	Created    time.Time
	Expiration int64
}

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanup           time.Duration
	items             map[string]Item
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanup:           cleanupInterval,
	}
	//есди интеравал >0 то запускаем очистку устаревших элементов
	if cleanupInterval > 0 {
		cache.startGC()
	}
	return &cache
}

//запихиваем в кеш значение, возможна перезапись существующего элемента
func (s *Cache) Set(key string, value []device.Device, duration time.Duration) {
	var expiration int64
	//определение продолжительности жизни айтема
	if duration == 0 {
		duration = s.defaultExpiration
	}
	// Время жизни кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	s.Lock()

	defer s.Unlock()

	s.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (s *Cache) GetAll() (items []Item) {
	for _, value := range s.items {
		items = append(items, value)
	}
	return
}

func (s *Cache) Get(key string) ([]device.Device, bool) {
	s.RLock()

	defer s.RUnlock()

	item, found := s.items[key]

	if !found {
		return nil, false
	}
	//проверяем бессрочно ли хранится элемент
	if item.Expiration > 0 {
		//проверяем актуальность элемента
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

func (s *Cache) Delete(key string) error {
	s.Lock()

	defer s.Unlock()

	if _, found := s.items[key]; !found {
		return errors.New("this item does not exist")
	}

	delete(s.items, key)

	return nil
}

func (s *Cache) startGC() {
	go s.GC()
}

func (s *Cache) GC() {
	for {
		<-time.After(s.cleanup)

		if s.items == nil {
			return
		}

		if keys := s.expiredKeys(); len(keys) != 0 {
			s.clear(keys)
		}
	}

}

func (s *Cache) expiredKeys() (keys []string) {
	s.RLock()

	defer s.RUnlock()
	for i, value := range s.items {
		if time.Now().UnixNano() > value.Expiration && value.Expiration > 0 {
			keys = append(keys, i)
		}
	}

	return
}

func (s *Cache) clear(keys []string) {
	s.Lock()

	defer s.Unlock()
	for _, key := range keys {
		delete(s.items, key)
	}
}
