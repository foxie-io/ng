package orders

import (
	"context"
	"errors"
	"example/advanced/components/orders/dtos"
	"example/advanced/dal"
	"example/advanced/models"

	"github.com/foxie-io/gormqs"
	nghttp "github.com/foxie-io/ng/http"
	"gorm.io/gorm"
)

type OrderService struct {
	orderList []*models.Order
	orderDao  *dal.OrderDao
}

func NewOrderService(orderDao *dal.OrderDao) *OrderService {
	return &OrderService{
		orderDao: orderDao,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req dtos.CreateOrderRequest) (*dtos.CreateOrderResponse, error) {
	record := &models.Order{
		UserID:   req.UserID,
		Product:  req.Product,
		Quantity: req.Quantity,
	}

	if err := s.orderDao.CreateOne(ctx, record); err != nil {
		return nil, err
	}

	return &dtos.CreateOrderResponse{
		ID:       record.ID,
		UserID:   record.UserID,
		Product:  record.Product,
		Quantity: record.Quantity,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (*dtos.GetOrderResponse, error) {
	var (
		record dtos.GetOrderResponse
	)

	err := s.orderDao.GetOneTo(ctx, &record, gormqs.WhereID(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nghttp.NewErrNotFound().Update(nghttp.Meta(
			"entity", "Order",
			"id", id,
		))
	}

	if err != nil {
		return nil, err
	}
	return &record, nil
}
func (s *OrderService) GetOrders(ctx context.Context, dto *dtos.ListOrdersRequest) *dtos.GetAllOrdersResponse {
	var (
		orders []dtos.GetOrderResponse
		limit  = dto.PageSize
		offset = (dto.Page - 1) * dto.PageSize
	)

	if err := s.orderDao.GetManyTo(ctx, &orders, gormqs.LimitAndOffset(limit, offset)); err != nil {
		return &dtos.GetAllOrdersResponse{Orders: []dtos.GetOrderResponse{}}
	}

	return &dtos.GetAllOrdersResponse{Orders: orders}
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, req *dtos.UpdateOrderRequest) (*dtos.UpdateOrderResponse, error) {
	_, err := s.orderDao.GetOne(ctx, gormqs.WhereID(id), gormqs.Select("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nghttp.NewErrNotFound().Update(nghttp.Meta(
			"entity", "Order",
			"id", id,
		))
	}
	if err != nil {
		return nil, err
	}

	updatedOrder := &models.Order{
		ID:       id,
		Product:  req.Product,
		Quantity: req.Quantity,
	}

	if _, err := s.orderDao.Update(ctx, updatedOrder,
		gormqs.WhereID(id),
		gormqs.Select("product", "quantity"),
	); err != nil {
		return nil, err
	}

	return &dtos.UpdateOrderResponse{
		ID:       updatedOrder.ID,
		Product:  updatedOrder.Product,
		Quantity: updatedOrder.Quantity,
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) (*dtos.DeleteOrderResponse, error) {
	_, err := s.orderDao.GetOne(ctx, gormqs.WhereID(id), gormqs.Select("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nghttp.NewErrNotFound().Update(nghttp.Meta(
			"entity", "Order",
			"id", id,
		))
	}

	if err != nil {
		return nil, err
	}

	if _, err := s.orderDao.Delete(ctx, gormqs.WhereID(id)); err != nil {
		return nil, err
	}

	return &dtos.DeleteOrderResponse{Success: true}, nil
}
