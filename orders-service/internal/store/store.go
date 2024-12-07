package store

type OrderStore interface {
	GetOrders()
	PlaceOrders()
	GetByID()
}
