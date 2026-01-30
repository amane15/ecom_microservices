package service

import (
	"context"

	"github.com/amane15/ecom_microservice/services/orders/internal/data"
)

type CreateOrderInput struct{}

type CreateOrderItemInput struct{}

type OrderRepository interface {
	Get(ctx context.Context, id int64) (*data.Order, error)
	Create(ctx context.Context, input *CreateOrderInput) (*data.Order, error)
	Delete(ctx context.Context, id int64) error
}

type ItemRepository interface {
	Get(ctx context.Context, id int64) (*data.OrderItem, error)
	Create(ctx context.Context, input *CreateOrderItemInput) (*data.OrderItem, error)
	Delete(ctx context.Context, id int64) error
}

type OrderService struct {
	orderRepository OrderRepository
	itemRepository  ItemRepository
}

func NewOrderService(orderRepository OrderRepository, itemRepository ItemRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		itemRepository:  itemRepository,
	}
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (*data.Order, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	order, err := s.orderRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
