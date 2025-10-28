Пошаговая инструкция по запуску:

1. cd D:\WB\order-service  (перейти местоположение файла)
2. go mod tidy (обновить зависимости)
3. nats-streaming-server.exe -store file -dir datastore -cluster_id test-cluster  (запуск NATS Streaming)
4. go run cmd/service/main.go (Запустить)