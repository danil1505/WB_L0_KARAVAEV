package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"order-service/config"
	"order-service/internal/models"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	connStr := cfg.GetConnectionString()
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return &Database{conn: conn}, nil
}

func (db *Database) SaveOrder(order *models.Order) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", order.OrderUID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Order %s already exists, skipping", order.OrderUID)
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
				sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *Database) GetAllOrders() ([]*models.Order, error) {
	rows, err := db.conn.Query(`
		SELECT order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		ORDER BY date_created DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		)
		if err != nil {
			return nil, err
		}

		err = db.conn.QueryRow(`
			SELECT name, phone, zip, city, address, region, email
			FROM delivery WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		err = db.conn.QueryRow(`
			SELECT transaction, request_id, currency, provider, amount,
				payment_dt, bank, delivery_cost, goods_total, custom_fee
			FROM payment WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		itemRows, err := db.conn.Query(`
			SELECT chrt_id, track_number, price, rid, name, sale, size,
				total_price, nm_id, brand, status
			FROM items WHERE order_uid = $1
		`, order.OrderUID)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {
			item := models.Item{}
			err := itemRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid,
				&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
				&item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		orders = append(orders, order)
	}

	return orders, nil
}

func (db *Database) Close() error {
	if db.conn != nil {
		log.Println("Disconnecting from PostgreSQL")
		return db.conn.Close()
	}
	return nil
}
