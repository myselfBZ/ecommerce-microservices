package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgreStore struct {
	db *sql.DB
}

func NewPostgre() (*PostgreStore, error) {
	db, err := sql.Open("postgres", "host=localhost port=32768 user=postgres password=new_password dbname=products sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &PostgreStore{
		db: db,
	}, nil
}

func (s *PostgreStore) GetOrders() ([]*Order, error) {
	q := `SELECT id, user_id, product_id, quantity, price, address FROM orders`
	r, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	var orders []*Order
	for r.Next() {
		var order Order
		err := r.Scan(&order.ID, &order.UserId, &order.ProductId, &order.ProductQuantity, &order.Address)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (s *PostgreStore) GetByID(id int) (*Order, error) {
	q := `SELECT id, user_id, product_id, quantity, price, address FROM orders WHERE id = $1`
	r := s.db.QueryRow(q, id)
	var order Order
	if err := r.Scan(&order.ID, &order.UserId, &order.ProductId, &order.ProductQuantity, &order.Address); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *PostgreStore) PlaceOrder(o *Order) error {
	q := `INSERT INTO orders(product_id, user_id, quantity, price, address) VALUES($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(q, o.ProductId, o.UserId, o.ProductQuantity, o.Address)
	return err
}
