package cache

import (
	"ex0/model"
	"ex0/storage"
	"log"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	c   *cache.Cache
	db  *storage.Storage
}

func NewCache(db *storage.Storage) *Cache {
	return &Cache{c: cache.New(cache.NoExpiration, cache.NoExpiration), db: db}
}

func (c *Cache) Init() {
	orders, err := c.db.FillCache()
	if err != nil {
		log.Fatalf("Failed to fill cache: %v", err)
	}

	c.loadAndCacheData(orders)
}

func (c *Cache) Load(o model.Model) {
	c.c.Set(o.OrderUid, o, cache.NoExpiration)
}

func (c *Cache) loadAndCacheData(orders []model.Model) {
	for _, order := range orders {
		c.c.Set(order.OrderUid, order, cache.NoExpiration)
	}
}

func (c *Cache) GetOrderFromCache(orderID string) (model.Model, bool) {
	order, found := c.c.Get(orderID)
	if !found {
		return model.Model{}, false
	}

	return order.(model.Model), true
}
