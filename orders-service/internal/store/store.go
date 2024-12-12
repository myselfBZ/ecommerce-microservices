package store

type Store interface {
	GetOrders() ([]*Order, error)
	PlaceOrder(*Order) error
	GetByID(int) (*Order, error)
}

type Order struct {
	ID              int
	ProductId       int
	UserId          int
	Price           float64
	ProductQuantity int
	Address         string
}
