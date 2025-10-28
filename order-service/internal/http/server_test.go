package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/config"
	"order-service/internal/models"
)

type mockCache struct {
	data map[string]*models.Order
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string]*models.Order),
	}
}

func (m *mockCache) Set(order *models.Order) {
	m.data[order.OrderUID] = order
}

func (m *mockCache) Get(orderUID string) (*models.Order, bool) {
	order, exists := m.data[orderUID]
	return order, exists
}

func (m *mockCache) GetAll() map[string]*models.Order {
	return m.data
}

func TestServerGetOrder(t *testing.T) {
	cache := newMockCache()
	order := &models.Order{
		OrderUID:    "TEST_ORDER",
		TrackNumber: "TRACK123",
		CustomerID:  "Test Customer",
	}
	cache.Set(order)

	cfg := &config.HTTPConfig{Port: "8080"}
	server := NewServer(cfg, cache)

	req := httptest.NewRequest("GET", "/api/orders/TEST_ORDER", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.Order
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal("Error decoding response:", err)
	}

	if response.OrderUID != "TEST_ORDER" {
		t.Errorf("Expected OrderUID TEST_ORDER, got %s", response.OrderUID)
	}
}

func TestServerGetOrderNotFound(t *testing.T) {
	cache := newMockCache()
	cfg := &config.HTTPConfig{Port: "8080"}
	server := NewServer(cfg, cache)

	req := httptest.NewRequest("GET", "/api/orders/NON_EXISTENT", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestServerGetAllOrders(t *testing.T) {
	cache := newMockCache()

	for i := 1; i <= 5; i++ {
		order := &models.Order{OrderUID: string(rune(i))}
		cache.Set(order)
	}

	cfg := &config.HTTPConfig{Port: "8080"}
	server := NewServer(cfg, cache)

	req := httptest.NewRequest("GET", "/api/orders", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal("Error decoding response:", err)
	}

	count := int(response["count"].(float64))
	if count != 5 {
		t.Errorf("Expected 5 orders, got %d", count)
	}
}

func TestServerIndexPage(t *testing.T) {
	cache := newMockCache()
	cfg := &config.HTTPConfig{Port: "8080"}
	server := NewServer(cfg, cache)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("Expected Content-Type: text/html; charset=utf-8")
	}
}
