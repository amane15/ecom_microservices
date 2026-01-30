package service

import (
	"context"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
)

type CreateCartInput struct {
	UserID int64 `json:"user_id"`
}

type CreateItemInput struct {
	CartID    int64  `json:"cart_id"`
	ProductID *int64 `json:"product_id"`
	VariantID *int64 `json:"variant_id"`
	Quantity  *int32 `json:"quantity"`
}

type CartRepository interface {
	Get(ctx context.Context, id int64) (*data.Cart, error)
	Create(ctx context.Context, input *CreateCartInput) (*data.Cart, error)
	Delete(ctx context.Context, id int64) error
}

type ItemRepository interface {
	Get(ctx context.Context, id int64) (*data.Item, error)
	Create(ctx context.Context, input *CreateItemInput) (*data.Item, error)
	Delete(ctx context.Context, id int64) error
}

type CartService struct {
	cartRepository CartRepository
	itemRepository ItemRepository
}

func NewCartService(cartRepository CartRepository, itemRepository ItemRepository) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		itemRepository: itemRepository,
	}
}

func (s *CartService) GetCart(ctx context.Context, id int64) (*data.Cart, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	cart, err := s.cartRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) CreateCart(ctx context.Context, input *CreateCartInput) (*data.Cart, error) {
	// verify user from grpc later
	cart, err := s.cartRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) DeleteCart(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.cartRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) GetItem(ctx context.Context, id int64) (*data.Item, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	item, err := s.itemRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *CartService) CreateItem(ctx context.Context, input *CreateItemInput) (*data.Item, error) {
	v := validator.New()

	// check for product_id and variant_id from grpc later
	if input.Quantity != nil {
		v.Check(*input.Quantity < 0, "quantity", "must be greater than zero")
	}

	if !v.Valid() {
		return nil, &data.ValidationError{
			Fields: v.Errors,
		}
	}

	item, err := s.itemRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *CartService) DeleteItem(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.itemRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
