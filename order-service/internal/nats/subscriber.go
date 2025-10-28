package nats

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/stan.go"
	"order-service/config"
	"order-service/internal/models"
)

type Cache interface {
	Set(order *models.Order)
}

type Database interface {
	SaveOrder(order *models.Order) error
}

type Subscriber struct {
	conn       stan.Conn
	cache      Cache
	db         Database
	subject    string
	queueGroup string
}

func NewSubscriber(cfg *config.NATSConfig, cache Cache, db Database) (*Subscriber, error) {
	conn, err := stan.Connect(
		cfg.ClusterID,
		cfg.ClientID,
		stan.NatsURL(cfg.URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Printf("Connection lost to NATS: %v", err)
		}),
		stan.Pings(10, 5),
	)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to NATS Streaming")

	return &Subscriber{
		conn:       conn,
		cache:      cache,
		db:         db,
		subject:    cfg.Subject,
		queueGroup: "order-service-group",
	}, nil
}

func (s *Subscriber) Subscribe() error {
	_, err := s.conn.QueueSubscribe(
		s.subject,
		s.queueGroup,
		s.handleMessage,
		stan.DurableName("order-service-durable"),
		stan.SetManualAckMode(),
		stan.AckWait(30*time.Second),
		stan.MaxInflight(25),
	)
	if err != nil {
		return err
	}

	log.Printf("Subscribed to channel: %s (group: %s)", s.subject, s.queueGroup)
	return nil
}

func (s *Subscriber) handleMessage(msg *stan.Msg) {
	log.Printf("Message received: seq=%d", msg.Sequence)

	var order models.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		log.Printf("JSON parsing error: %v", err)
		msg.Ack()
		return
	}

	if err := s.db.SaveOrder(&order); err != nil {
		log.Printf("Database save error: %v", err)
		return
	}

	s.cache.Set(&order)

	log.Printf("Order %s saved successfully", order.OrderUID)
	msg.Ack()
}

func (s *Subscriber) Close() {
	if s.conn != nil {
		s.conn.Close()
		log.Println("Disconnected from NATS Streaming")
	}
}
