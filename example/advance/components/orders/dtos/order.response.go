package dtos

import "github.com/foxie-io/gormqs"

type CreateOrderResponse struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userId"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

var _ interface{ gormqs.Model } = (*GetOrderResponse)(nil)

type GetOrderResponse struct {
	ID       int    `json:"id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func (r GetOrderResponse) TableName() string {
	return "orders"
}

type GetAllOrdersResponse struct {
	Orders []GetOrderResponse `json:"orders"`
}

type UpdateOrderResponse struct {
	ID       int    `json:"id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type DeleteOrderResponse struct {
	Success bool `json:"success"`
}
