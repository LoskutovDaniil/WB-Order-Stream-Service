package storage

import (
	"database/sql"
	"ex0/model"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() {
	s.db.Close()
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddOrder(order model.Model) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	deliveryQuery := `INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var deliveryID int
	err = tx.QueryRow(deliveryQuery, order.Dev.Name, order.Dev.Phone, order.Dev.Zip, order.Dev.City, order.Dev.Address, order.Dev.Region, order.Dev.Email).Scan(&deliveryID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert delivery: %v", err)
	}

	paymentQuery := `INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var paymentID int
	err = tx.QueryRow(paymentQuery, order.Pay.Transaction, order.Pay.RequestId, order.Pay.Currency, order.Pay.Provider, order.Pay.Amount, order.Pay.PaymentDt, order.Pay.Bank, order.Pay.DeliveryCost, order.Pay.GoodsTotal, order.Pay.CustomFee).Scan(&paymentID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert payment: %v", err)
	}

	orderQuery := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = tx.Exec(orderQuery, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard, deliveryID, paymentID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %v", err)
	}

	itemQuery := `INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range order.It {
		_, err = tx.Exec(itemQuery, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status, order.OrderUid)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert item: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *Storage) GetOrder(id string) (model.Model, error) {
	var order model.Model
	var order_id int

	orderQuery := `
		SELECT o.id,o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		       d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		       p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_id = p.id
		WHERE o.order_uid = $1
	`

	row := s.db.QueryRow(orderQuery, id)

	err := row.Scan(
		&order_id, &order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard,
		&order.Dev.Name, &order.Dev.Phone, &order.Dev.Zip, &order.Dev.City, &order.Dev.Address, &order.Dev.Region, &order.Dev.Email,
		&order.Pay.Transaction, &order.Pay.RequestId, &order.Pay.Currency, &order.Pay.Provider, &order.Pay.Amount, &order.Pay.PaymentDt, &order.Pay.Bank, &order.Pay.DeliveryCost, &order.Pay.GoodsTotal, &order.Pay.CustomFee,
	)
	if err != nil {
		return order, fmt.Errorf("failed to get order: %v", err)
	}

	itemsQuery := `
    SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
    FROM items
    WHERE order_uid = $1
`

	rows, err := s.db.Query(itemsQuery, order_id)
	if err != nil {
		return order, fmt.Errorf("failed to get items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Items
		if err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status); err != nil {
			return order, fmt.Errorf("failed to scan item: %v", err)
		}
		order.It = append(order.It, item)
	}

	return order, nil
}

func (s *Storage) FillCache() ([]model.Model, error) {
	var orders []model.Model

	rows, err := s.db.Query("SELECT order_uid FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderUid string
		if err := rows.Scan(&orderUid); err != nil {
			return nil, err
		}

		order, err := s.GetOrder(orderUid)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
