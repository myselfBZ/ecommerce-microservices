package store

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

const (
	UniqueError = "23505"
)

var fieldColums = map[string]string{
	"Name":     "first_name",
	"LastName": "last_name",
	"Email":    "email",
}

type PostgreStore struct {
	db *sql.DB
}

func NewPostgre() (*PostgreStore, error) {
	db, err := sql.Open("postgres", "host=localhost port=32768 user=postgres password=new_password dbname=users sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &PostgreStore{
		db: db,
	}, nil
}

func (s *PostgreStore) Create(u *User) error {
	q := `INSERT INTO users(first_name, last_name, email, password) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(q, u.Name, u.LastName, u.Email, u.Password)
	return err
}

func (s *PostgreStore) Delete(id int) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(q, id)
	return err
}

func (s *PostgreStore) GetById(id int) (*User, error) {
	q := "SELECT id, first_name, last_name, email FROM users;"
	r := s.db.QueryRow(q)
	var u User
	err := r.Scan(&u.ID, &u.Name, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// updating method is more difficult to implement than other methods.
// we have some reflective programming going on here
func (s *PostgreStore) Update(u *User, id int) error {
	template := `UPDATE users SET %s WHERE id = %d`
	nonNilFields := eliminateNil(u)
	colums, fields := adjustQuery(nonNilFields)
	if colums == "" {
		return nil
	}
	query := fmt.Sprintf(template, colums, id)

	_, err := s.db.Exec(query, fields...)
	return err
}

func (s *PostgreStore) GetByEmail(e string) (*User, error) {
	q := "SELECT password FROM users WHERE email = $1"
	r, err := s.db.Query(q, e)
	if err != nil {
		return nil, err
	}
	var u User

	if r.Next() {
		err = r.Scan(&u.Password)
		if err != nil {
			return nil, err
		}
	}

	return &u, nil
}

// for update method.

func eliminateNil(b *User) map[string]interface{} {
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
