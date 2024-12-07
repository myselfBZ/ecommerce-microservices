package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func DBConnect() *sql.DB {
	psql, err := sql.Open("postgres", "host=localhost port=32768 user=postgres password=new_password dbname=users sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	log.Println("connected")
	return psql
}

func createUsersTable() error {
	db := DBConnect()
	q := `CREATE TABLE IF NOT EXISTS users (
		ID SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		role 	INT,
		email  VARCHAR(255),
		password VARCHAR(255)
	)`

	_, err := db.Exec(q)
	return err
}

func main() {
	err := createUsersTable()
	if err != nil {
		log.Fatal("error creating users table: ", err)
	}
	log.Printf("let's go")
}
