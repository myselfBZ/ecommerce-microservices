package store

import "database/sql"

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
