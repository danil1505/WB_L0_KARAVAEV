package cache

import (
	"testing"
	"time"

	"order-service/internal/models"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := NewCache()

	order := &models.Order{
		OrderUID:    "TEST_001",
		TrackNumber: "TRACK001",
		CustomerID:  "Test Client",
	}

	cache.Set(order)

	retrieved, exists := cache.Get("TEST_001")
	if !exists {
		t.Fatal("Order should exist in cache")
	}

	if retrieved.OrderUID != "TEST_001" {
		t.Errorf("Expected OrderUID TEST_001, got %s", retrieved.OrderUID)
	}
}

func TestCacheGetNonExistent(t *testing.T) {
	cache := NewCache()

	_, exists := cache.Get("NON_EXISTENT")
	if exists {
		t.Error("Non-existent order should not be returned")
	}
}

func TestCacheGetAll(t *testing.T) {
	cache := NewCache()

	orders := []*models.Order{
		{OrderUID: "ORDER_1"},
		{OrderUID: "ORDER_2"},
		{OrderUID: "ORDER_3"},
	}

	for _, order := range orders {
		cache.Set(order)
	}

	all := cache.GetAll()
	if len(all) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(all))
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache()

	for i := 0; i < 100; i++ {
		go func(id int) {
			order := &models.Order{
				OrderUID: string(rune(id)),
			}
			cache.Set(order)
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(id int) {
			cache.Get(string(rune(id)))
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache()
	order := &models.Order{OrderUID: "BENCH_TEST"}
	cache.Set(order)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("BENCH_TEST")
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		order := &models.Order{OrderUID: string(rune(i))}
		cache.Set(order)
	}
}
