


Write-Host "Stress-Test" -ForegroundColor Green

Write-Host "`n[1] Test GET /api/orders (100 req/sec, 30 сек)" -ForegroundColor Yellow
echo "GET http://localhost:8080/api/orders" | vegeta attack -duration=30s -rate=100 | vegeta report

Write-Host "`n[2] Test GET /api/orders/ЗАКАЗ_01 (200 req/sec, 30 сек)" -ForegroundColor Yellow
echo "GET http://localhost:8080/api/orders/ЗАКАЗ_01" | vegeta attack -duration=30s -rate=200 | vegeta report

Write-Host "`n[3] Test GET / (150 req/sec, 30 сек)" -ForegroundColor Yellow
echo "GET http://localhost:8080/" | vegeta attack -duration=30s -rate=150 | vegeta report

Write-Host "`n[4] Extreme test (500 req/sec, 10 сек)" -ForegroundColor Yellow
echo "GET http://localhost:8080/api/orders" | vegeta attack -duration=10s -rate=500 | vegeta report

Write-Host "`nTests completed" -ForegroundColor Green