package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"order-service/config"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/http"
	"order-service/internal/nats"
)

func main() {
	log.Println("Запуск Order Service...")

	cfg := config.GetConfig()

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	orderCache := cache.NewCache()
	if err := orderCache.RestoreFromDB(db); err != nil {
		log.Printf("Предупреждение при восстановлении кэша: %v", err)
	}

	subscriber, err := nats.NewSubscriber(&cfg.NATS, orderCache, db)
	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}
	defer subscriber.Close()

	if err := subscriber.Subscribe(); err != nil {
		log.Fatal("Ошибка подписки:", err)
	}

	server := http.NewServer(&cfg.HTTP, orderCache)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Ошибка HTTP сервера:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Остановка сервиса...")
}
