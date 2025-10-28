package cache

import (
	"log"
	"sync"

	"order-service/internal/models"
)

type Database interface {
	GetAllOrders() ([]*models.Order, error)
}

type Cache struct {
	data map[string]*models.Order
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*models.Order),
	}
}

func (c *Cache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func (c *Cache) Get(orderUID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.data[orderUID]
	return order, exists
}

func (c *Cache) GetAll() map[string]*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*models.Order, len(c.data))
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

func (c *Cache) RestoreFromDB(db Database) error {
	log.Println("Восстановление кэша из БД...")

	orders, err := db.GetAllOrders()
	if err != nil {
		return err
	}

	c.mu.Lock()
	for _, order := range orders {
		c.data[order.OrderUID] = order
	}
	c.mu.Unlock()

	log.Printf("Кэш восстановлен: загружено %d заказов", len(orders))
	return nil
}
