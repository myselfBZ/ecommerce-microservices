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

func (s *PostgreStore) GetOrders() {

}

func (s *PostgreStore) GetByID() {

}

func (s *PostgreStore) PlaceOrder(o *Order) error {
	q := `INSERT INTO orders(product_id, user_id, quantity, price, address) VALUES($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(q, o.ProductId, o.UserId, o.ProductQuantity, o.Address)
	return err
}
