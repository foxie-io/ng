package models

type Order struct {
	ID       int
	UserID   int
	Product  string
	Quantity int
}

func (o Order) TableName() string {
	return "orders"
}
