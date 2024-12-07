package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func connectToDB() {
	db, err := sql.Open("postgres", "host=localhost port=32768 user=postgres password=new_password dbname=products sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database")
	}
	DB = db
}

func init() {
	connectToDB()
}

func createProductsTable() error {
	q := `CREATE TABLE IF NOT EXISTS products(
		ID SERIAL PRIMARY KEY,
		name VARCHAR(255),
		description TEXT,
		price DECIMAL,
		quantity INT
	)`
	_, err := DB.Exec(q)
	return err
}

func main() {
	log.Println("error creating products table: ", createProductsTable())
}
