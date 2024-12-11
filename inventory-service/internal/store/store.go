package store

import "time"

type Product struct {
	ID          int
	Price       float64
	Description string
	Name        string
	Quantity    int
}

type StockTransaction struct {
	ID             int
	Product_id     int
	QuantityChange int
	Reason         string
	CreatedAt      time.Time
}

type Store interface {
	GetProductByID(int) (*Product, error)
	GetProducts() ([]*Product, error)
	CreateProduct(*Product) error
	CreateStockTransaction(*StockTransaction) error
	DeleteProduct(int) error
	UpdateProduct(*Product, int) error
}
