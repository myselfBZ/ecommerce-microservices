package store

type Store interface {
	GetOrders()
	PlaceOrder(*Order) error
	GetByID()
}

type Order struct {
	ID              int
	ProductId       int
	UserId          int
	Price           float64
	ProductQuantity int
	Address         string
}
