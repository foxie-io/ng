package dtos

type CreateOrderRequest struct {
	UserID   int    `json:"userId" validate:"required,gt=0"`
	Product  string `json:"product" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderRequest struct {
	Product  string `json:"product" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

type PathID struct {
	ID int `json:"id" param:"id" validate:"required,gt=0"`
}

type ListOrdersRequest struct {
	Page     int `query:"page" default:"1" validate:"gte=1"`
	PageSize int `query:"pageSize" default:"10" validate:"gte=1,lte=100"`
}
