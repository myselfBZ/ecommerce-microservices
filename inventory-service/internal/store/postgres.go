package store

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

type PostgreStore struct {
	db *sql.DB
}

var fieldColums = map[string]string{
	"Price":       "price",
	"Description": "description",
	"Name":        "name",
	"Quantity":    "quantity",
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

func (s *PostgreStore) GetProductByID(id int) (*Product, error) {
	q := `SELECT id, name, description, price, quantity FROM products WHERE id = $1`
	r := s.db.QueryRow(q, id)
	var p Product
	err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgreStore) CreateProduct(p *Product) error {
	q := `INSERT INTO products(name, description, price, quantity) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(q, p.Name, p.Description, p.Price, p.Quantity)
	return err
}

func (s *PostgreStore) DeleteProduct(id int) error {
	q := `DELETE FROM products WHERE id = $1`
	_, err := s.db.Exec(q, id)
	return err
}

func (s *PostgreStore) UpdateProduct(p *Product, id int) error {
	template := `UPDATE products SET %s WHERE id = %d`
	nonNilFields := eliminateNil(p)
	colums, fields := adjustQuery(nonNilFields)
	if colums == "" {
		return nil
	}
	query := fmt.Sprintf(template, colums, id)

	r := s.db.QueryRow(query, fields...)
	return r.Err()
}

func (s *PostgreStore) CreateStockTransaction(st *StockTransaction) error {
	q := `INSERT INTO transactions (product_id, quantity_change, reason, created_at) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(q, st.Product_id, st.QuantityChange, st.Reason, st.CreatedAt)
	return err
}

func eliminateNil(b *Product) map[string]interface{} {
	m := make(map[string]interface{})
	// doesn't work with pointers, HAS TO BE DEREFRENCED
	t := reflect.TypeOf(*b)
	v := reflect.ValueOf(*b)
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		if !f.IsZero() {
			m[t.Field(i).Name] = f.Interface()
		}
	}
	return m
}

// for update method
func adjustQuery(m map[string]interface{}) (string, []interface{}) {
	//columns
	columns := ""
	counter := 0
	fields := []interface{}{}
	for k, v := range m {
		counter++
		columnName := fieldColums[k]
		fields = append(fields, v)
		columns += fmt.Sprintf("%s = $%d, ", columnName, counter)
	}

	if len(columns) > 0 {
		return columns[:len(columns)-2], fields
	}

	return "", nil
}
