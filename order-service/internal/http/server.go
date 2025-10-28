package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"order-service/config"
	"order-service/internal/models"
)

type Cache interface {
	Get(orderUID string) (*models.Order, bool)
	GetAll() map[string]*models.Order
}

type Server struct {
	router *mux.Router
	cache  Cache
	port   string
}

func NewServer(cfg *config.HTTPConfig, cache Cache) *Server {
	server := &Server{
		router: mux.NewRouter(),
		cache:  cache,
		port:   cfg.Port,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.router.HandleFunc("/api/orders/{id}", s.handleGetOrder).Methods("GET")
	s.router.HandleFunc("/api/orders", s.handleGetAllOrders).Methods("GET")
	s.router.HandleFunc("/orders/{id}", s.handleOrderPage).Methods("GET")
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Задание L0 WILDBERRIES</title>
    <style>
        * { 
            margin: 0; 
            padding: 0; 
            box-sizing: border-box; 
        }
        body {
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #8b4a8f 0%, #6b3a6e 100%);
            min-height: 100vh;
            padding: 40px 20px;
            display: flex;
            justify-content: center;
            align-items: flex-start;
            position: relative;
            overflow-x: hidden;
        }
        body::before {
            content: 'L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB';
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            font-size: 80px;
            font-weight: bold;
            color: rgba(255, 255, 255, 0.03);
            line-height: 1.8;
            word-spacing: 80px;
            transform: rotate(-15deg);
            pointer-events: none;
            white-space: pre-wrap;
            overflow: hidden;
            z-index: 0;
        }
        .container {
            max-width: 1100px;
            width: 100%;
            background: #f5f5f5;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.4);
            position: relative;
            z-index: 1;
        }
        h1 { 
            color: #c73659;
            margin-bottom: 30px;
            text-align: center;
            font-size: 2em;
            font-weight: bold;
        }
        .search-box { 
            display: flex; 
            gap: 15px; 
            margin-bottom: 30px;
            justify-content: center;
            align-items: center;
        }
        input {
            padding: 15px 20px;
            border: 2px solid #ddd;
            border-radius: 8px;
            font-size: 15px;
            width: 400px;
            background: white;
        }
        input:focus {
            outline: none;
            border-color: #c73659;
        }
        .search-btn {
            padding: 15px 35px;
            background: linear-gradient(135deg, #d666a0 0%, #c73659 100%);
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 15px;
            font-weight: bold;
            transition: all 0.2s;
        }
        .search-btn:hover { 
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(199, 54, 89, 0.4);
        }
        .controls-section {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 20px;
            margin-bottom: 25px;
        }
        .sort-group {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        .sort-label {
            font-weight: bold;
            color: #333;
            font-size: 16px;
        }
        .filters {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }
        .filter-btn {
            padding: 12px 20px;
            background: white;
            color: #c73659;
            border: 2px solid #c73659;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
            font-weight: bold;
            transition: all 0.2s;
        }
        .filter-btn:hover {
            background: #c73659;
            color: white;
        }
        .filter-btn.active {
            background: #c73659;
            color: white;
        }
        .stat-card {
            background: linear-gradient(135deg, #d666a0 0%, #c73659 100%);
            color: white;
            padding: 20px 30px;
            border-radius: 15px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.2);
            min-width: 140px;
        }
        .stat-value { 
            font-size: 2.5em; 
            font-weight: bold; 
            line-height: 1;
            margin-bottom: 8px;
        }
        .stat-label {
            font-size: 0.9em;
        }
        .divider {
            height: 4px;
            background: linear-gradient(90deg, #c73659 0%, #d666a0 100%);
            margin: 25px 0;
            border-radius: 2px;
        }
        .section-header {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;
            margin-bottom: 25px;
        }
        .section-title {
            font-size: 1.4em;
            color: #555;
            font-weight: bold;
        }
        .orders-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 20px;
            margin-bottom: 25px;
        }
        @media (max-width: 900px) {
            .orders-grid {
                grid-template-columns: repeat(2, 1fr);
            }
        }
        @media (max-width: 600px) {
            .orders-grid {
                grid-template-columns: 1fr;
            }
        }
        .order-card {
            background: white;
            padding: 20px;
            border-radius: 12px;
            cursor: pointer;
            border: 2px solid #e0e0e0;
            transition: all 0.3s;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .order-card:hover {
            border-color: #c73659;
            transform: translateY(-5px);
            box-shadow: 0 8px 25px rgba(199, 54, 89, 0.3);
        }
        .order-header {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 12px;
        }
        .order-id {
            font-weight: bold;
            color: #c73659;
            font-size: 1.15em;
            word-break: break-all;
        }
        .order-info {
            color: #555;
            font-size: 0.95em;
            margin: 6px 0;
            line-height: 1.5;
        }
        .order-label {
            font-weight: 600;
            color: #333;
        }
        .show-more-btn {
            display: none;
            margin: 25px auto;
            padding: 15px 40px;
            background: linear-gradient(135deg, #d666a0 0%, #c73659 100%);
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 15px;
            font-weight: bold;
        }
        .show-more-btn.visible {
            display: block;
        }
        .show-more-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(199, 54, 89, 0.4);
        }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }
        .empty-icon {
            font-size: 4em;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Задание L0 WILDBERRIES</h1>
        
        <div class="search-box">
            <input type="text" id="orderIdInput" placeholder="Введите ID заказа для поиска">
            <button class="search-btn" onclick="searchOrder()">Поиск</button>
        </div>

        <div class="controls-section">
            <div class="sort-group">
                <span class="sort-label">Сортировка:</span>
                <div class="filters">
                    <button class="filter-btn active" onclick="setSortBy('date-desc')">По дате  ↓</button>
                    <button class="filter-btn" onclick="setSortBy('date-asc')">По дате ↑</button>
                    <button class="filter-btn" onclick="setSortBy('name-asc')">По имени А-Я</button>
                    <button class="filter-btn" onclick="setSortBy('name-desc')">По имени Я-А</button>
                    <button class="filter-btn" onclick="setSortBy('amount-desc')">По сумме ↓</button>
                    <button class="filter-btn" onclick="setSortBy('amount-asc')">По сумме ↑</button>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-value" id="cacheCount">0</div>
                <div class="stat-label">Заказов в кэше</div>
            </div>
        </div>

        <div class="divider"></div>

        <div class="section-header">
            <h2 class="section-title">Список заказов</h2>
        </div>

        <div id="ordersGrid" class="orders-grid"></div>
        
        <button id="showMoreBtn" class="show-more-btn" onclick="toggleOrders()">
            Показать все заказы
        </button>
        
        <div id="emptyState" class="empty-state" style="display: none;">
            <h3>Заказы не найдены</h3>
            <p>Ожидаем поступление заказов через NATS Streaming</p>
        </div>
    </div>

    <script>
        let allOrders = [];
        let showingAll = false;
        let currentSort = 'date-desc';
        const INITIAL_SHOW = 9;

        function searchOrder() {
            const id = document.getElementById('orderIdInput').value.trim();
            if (id) {
                window.location.href = '/orders/' + id;
            } else {
                alert('Введите ID заказа');
            }
        }

        function setSortBy(sortType) {
            currentSort = sortType;
            
            document.querySelectorAll('.filter-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            event.target.classList.add('active');
            
            displayOrders();
        }

        function sortOrders(orders) {
            const sorted = [...orders];
            
            switch(currentSort) {
                case 'date-desc':
                    sorted.sort((a, b) => new Date(b.date_created) - new Date(a.date_created));
                    break;
                case 'date-asc':
                    sorted.sort((a, b) => new Date(a.date_created) - new Date(b.date_created));
                    break;
                case 'name-asc':
                    sorted.sort((a, b) => (a.customer_id || '').localeCompare(b.customer_id || '', 'ru'));
                    break;
                case 'name-desc':
                    sorted.sort((a, b) => (b.customer_id || '').localeCompare(a.customer_id || '', 'ru'));
                    break;
                case 'amount-desc':
                    sorted.sort((a, b) => (b.payment?.amount || 0) - (a.payment?.amount || 0));
                    break;
                case 'amount-asc':
                    sorted.sort((a, b) => (a.payment?.amount || 0) - (b.payment?.amount || 0));
                    break;
            }
            
            return sorted;
        }

        async function loadAllOrders() {
            try {
                const response = await fetch('/api/orders');
                const data = await response.json();
                
                document.getElementById('cacheCount').textContent = data.count;
                
                allOrders = data.orders || [];
                displayOrders();
                
            } catch (err) {
                console.error('Error loading orders:', err);
            }
        }

        function displayOrders() {
            const grid = document.getElementById('ordersGrid');
            const emptyState = document.getElementById('emptyState');
            const showMoreBtn = document.getElementById('showMoreBtn');
            
            grid.innerHTML = '';
            
            if (allOrders.length === 0) {
                emptyState.style.display = 'block';
                showMoreBtn.classList.remove('visible');
                return;
            }
            
            emptyState.style.display = 'none';
            
            const sortedOrders = sortOrders(allOrders);
            const ordersToShow = showingAll ? sortedOrders : sortedOrders.slice(0, INITIAL_SHOW);
            
            ordersToShow.forEach(order => {
                const card = document.createElement('div');
                card.className = 'order-card';
                card.onclick = () => window.location.href = '/orders/' + order.order_uid;
                
                const date = new Date(order.date_created);
                const formattedDate = date.toLocaleDateString('ru-RU', { 
                    day: '2-digit', 
                    month: '2-digit', 
                    year: 'numeric' 
                });
                const amount = order.payment?.amount || 0;
                
                card.innerHTML = 
                    '<div class="order-header">' +
                    '<div class="order-id">' + order.order_uid + '</div>' +
                    '</div>' +
                    '<div class="order-info">' +
                    '<span class="order-label">Клиент:</span> ' + (order.customer_id || 'Не указан') +
                    '</div>' +
                    '<div class="order-info">' +
                    '<span class="order-label">Трек-номер:</span> ' + order.track_number +
                    '</div>' +
                    '<div class="order-info">' +
                    '<span class="order-label">Дата:</span> ' + formattedDate +
                    '</div>' +
                    '<div class="order-info">' +
                    '<span class="order-label">Сумма:</span> ' + amount + ' руб.' +
                    '</div>';
                
                grid.appendChild(card);
            });
            
            if (allOrders.length > INITIAL_SHOW) {
                showMoreBtn.classList.add('visible');
                const remaining = allOrders.length - INITIAL_SHOW;
                showMoreBtn.textContent = showingAll ? 
                    'Скрыть' : 
                    'Показать все (' + remaining + ' ещё)';
            } else {
                showMoreBtn.classList.remove('visible');
            }
        }

        function toggleOrders() {
            showingAll = !showingAll;
            displayOrders();
        }

        loadAllOrders();
        setInterval(loadAllOrders, 5000);
    </script>
</body>
</html>`)
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order, exists := s.cache.Get(vars["id"])
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (s *Server) handleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := s.cache.GetAll()
	list := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		list = append(list, order)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(list),
		"orders": list,
	})
}

func (s *Server) handleOrderPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order, exists := s.cache.Get(vars["id"])
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.New("order").Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Заказ {{.OrderUID}} - L0 WILDBERRIES</title>
<style>
* { 
    margin: 0; 
    padding: 0; 
    box-sizing: border-box; 
}
body {
    font-family: Arial, sans-serif;
    background: linear-gradient(135deg, #8b4a8f 0%, #6b3a6e 100%);
    min-height: 100vh;
    padding: 40px 20px;
    position: relative;
    overflow-x: hidden;
}
body::before {
    content: 'L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB L0 WB';
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    font-size: 80px;
    font-weight: bold;
    color: rgba(255, 255, 255, 0.03);
    line-height: 1.8;
    word-spacing: 80px;
    transform: rotate(-15deg);
    pointer-events: none;
    white-space: pre-wrap;
    overflow: hidden;
    z-index: 0;
}
.container { 
    max-width: 1000px; 
    margin: 0 auto; 
    background: #f5f5f5; 
    padding: 40px; 
    border-radius: 20px;
    box-shadow: 0 20px 60px rgba(0,0,0,0.4);
    position: relative;
    z-index: 1;
}
h1 { 
    color: #c73659; 
    margin-bottom: 30px;
    text-align: center;
    font-size: 1.8em;
}
.back-btn {
    display: inline-block;
    padding: 12px 25px;
    background: linear-gradient(135deg, #d666a0 0%, #c73659 100%);
    color: white;
    text-decoration: none;
    border-radius: 8px;
    margin-bottom: 25px;
    font-weight: bold;
    transition: all 0.2s;
}
.back-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(199, 54, 89, 0.4);
}
.divider {
    height: 4px;
    background: linear-gradient(90deg, #c73659 0%, #d666a0 100%);
    margin: 25px 0;
    border-radius: 2px;
}
.section { 
    background: white; 
    padding: 25px; 
    margin: 20px 0; 
    border-radius: 12px;
    border: 2px solid #e0e0e0;
}
h2 { 
    color: #c73659; 
    margin-bottom: 20px;
    font-size: 1.3em;
    display: flex;
    align-items: center;
    gap: 10px;
}
.info-grid { 
    display: grid; 
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); 
    gap: 15px; 
}
.info-item { 
    background: #f8f9fa; 
    padding: 15px; 
    border-radius: 8px;
    border-left: 4px solid #c73659;
}
.label { 
    color: #666; 
    font-size: 0.85em;
    font-weight: 600;
    text-transform: uppercase;
    margin-bottom: 5px;
}
.value { 
    color: #333; 
    font-weight: bold;
    font-size: 1.05em;
}
.item { 
    background: #f8f9fa; 
    padding: 18px; 
    margin: 12px 0; 
    border-radius: 8px;
    border-left: 4px solid #d666a0;
}
.item-header {
    font-weight: bold;
    color: #c73659;
    font-size: 1.1em;
    margin-bottom: 8px;
}
.item-details {
    color: #555;
    line-height: 1.6;
}
</style>
</head>
<body>
<div class="container">
<a href="/" class="back-btn">Назад к списку</a>
<h1>Заказ {{.OrderUID}}</h1>

<div class="divider"></div>

<div class="section">
<h2>Основная информация</h2>
<div class="info-grid">
<div class="info-item">
    <div class="label">Order UID</div>
    <div class="value">{{.OrderUID}}</div>
</div>
<div class="info-item">
    <div class="label">Трек-номер</div>
    <div class="value">{{.TrackNumber}}</div>
</div>
<div class="info-item">
    <div class="label">Клиент</div>
    <div class="value">{{.CustomerID}}</div>
</div>
<div class="info-item">
    <div class="label">Служба доставки</div>
    <div class="value">{{.DeliveryService}}</div>
</div>
</div>
</div>

<div class="section">
<h2>Доставка</h2>
<div class="info-grid">
<div class="info-item">
    <div class="label">Получатель</div>
    <div class="value">{{.Delivery.Name}}</div>
</div>

<div class="info-item">
    <div class="label">Телефон</div>
    <div class="value">{{.Delivery.Phone}}</div>
</div>

<div class="info-item">
    <div class="label">Email</div>
    <div class="value">{{.Delivery.Email}}</div>
</div>

<div class="info-item">
    <div class="label">Город</div>
    <div class="value">{{.Delivery.City}}</div>
</div>

<div class="info-item">
    <div class="label">Адрес</div>
    <div class="value">{{.Delivery.Address}}</div>
</div>

<div class="info-item">
    <div class="label">Регион</div>
    <div class="value">{{.Delivery.Region}}</div>
</div>

</div>
</div>

<div class="section">
<h2>Оплата</h2>
<div class="info-grid">
<div class="info-item">
    <div class="label">Сумма</div>
    <div class="value">{{.Payment.Amount}} {{.Payment.Currency}}</div>
</div>

<div class="info-item">
    <div class="label">Провайдер</div>
    <div class="value">{{.Payment.Provider}}</div>
</div>

<div class="info-item">
    <div class="label">Банк</div>
    <div class="value">{{.Payment.Bank}}</div>
</div>

<div class="info-item">
    <div class="label">Стоимость доставки</div>
    <div class="value">{{.Payment.DeliveryCost}} {{.Payment.Currency}}</div>
</div>

<div class="info-item">
    <div class="label">Товары</div>
    <div class="value">{{.Payment.GoodsTotal}} {{.Payment.Currency}}</div>
</div>

<div class="info-item">
    <div class="label">ID транзакции</div>
    <div class="value">{{.Payment.Transaction}}</div>
</div>

</div>
</div>

<div class="section">
<h2>Товары ({{len .Items}} шт.)</h2>
{{range .Items}}
<div class="item">
<div class="item-header">{{.Name}} - {{.Brand}}</div>
<div class="item-details">
    <strong>Цена:</strong> {{.Price}} руб. | 
    <strong>Скидка:</strong> {{.Sale}}% | 
    <strong>Итого:</strong> {{.TotalPrice}} руб.<br>
    <strong>Размер:</strong> {{.Size}} | 
    <strong>Трек:</strong> {{.TrackNumber}}
</div>

</div>
{{end}}
</div>

</div>
</body>
</html>`))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, order)
}

func (s *Server) Start() error {
	log.Printf("HTTP server started on http://localhost:%s", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
