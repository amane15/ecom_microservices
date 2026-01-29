package service

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/shopspring/decimal"
)

type UpdateProductInput struct {
	Name             *string             `json:"name"`
	Description      *string             `json:"description"`
	ShortDescription *string             `json:"short_description"`
	MetaTitle        *string             `json:"meta_title"`
	MetaDescription  *string             `json:"meta_description"`
	Status           *data.ProductStatus `json:"status"`
	DefaultVariantID *int64              `json:"default_variant_id"`
}

type CreateProductInput struct {
	Name             string              `json:"name"`
	Slug             string              `json:"slug"`
	Description      *string             `json:"description"`
	ShortDescription *string             `json:"short_description"`
	MetaTitle        *string             `json:"meta_title"`
	MetaDescription  *string             `json:"meta_description"`
	Status           *data.ProductStatus `json:"status"`
	DefaultVariantID *int64              `json:"default_variant_id"`
}

type UpdateStatusInput struct {
	Status *data.ProductStatus `json:"status"`
}

type SetDefaultVariantInput struct {
	VariantID int64 `json:"variant_id"`
}

type UpdateVariantInput struct {
	Name     *string          `json:"name"`
	Price    *decimal.Decimal `json:"price"`
	IsActive *bool            `json:"is_active"`
}

type CreateVariantInput struct {
	ProductID int64           `json:"product_id"`
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	IsActive  *bool           `json:"is_active"`
}

type VariantRepository interface {
	Get(ctx context.Context, id int64) (*data.Variant, error)
	Create(ctx context.Context, input *CreateVariantInput) (*data.Variant, error)
	Update(ctx context.Context, id int64, input *UpdateVariantInput) (*data.Variant, error)
	Delete(ctx context.Context, id int64) error
	ListByProduct(ctx context.Context, productID int64) ([]*data.Variant, error)
}

type ProductRepository interface {
	Get(ctx context.Context, id int64) (*data.Product, error)
	Create(ctx context.Context, product *CreateProductInput) (*data.Product, error)
	Update(ctx context.Context, id int64, input *UpdateProductInput) (*data.Product, error)
	List(ctx context.Context, limit, offset int32) ([]*data.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductService struct {
	productRepository ProductRepository
	variantRepository VariantRepository
}

func NewProductService(r ProductRepository, v VariantRepository) *ProductService {
	return &ProductService{
		productRepository: r,
		variantRepository: v,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*data.Product, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	product, err := s.productRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, input *CreateProductInput) (*data.Product, error) {
	v := validator.New()

	validateName(v, input.Name)
	validateSlug(v, input.Slug)
	if input.Status != nil {
		validateStatus(v, *input.Status)
	}

	if !v.Valid() {
		return nil, &data.ValidationErrors{
			Fields: v.Errors,
		}
	}

	product, err := s.productRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, input *UpdateProductInput) (*data.Product, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	v := validator.New()

	if input.Name != nil {
		validateName(v, *input.Name)
	}

	if input.DefaultVariantID != nil {
		// verify varaint_id with variant service later
	}

	if input.Status != nil {
		status := *input.Status
		switch status {
		case data.ProductStatusDraft, data.ProductStatusArchived, data.ProductStatusActive:
			if status == data.ProductStatusActive {
				count, err := 0, error(nil)
				if err != nil {
					switch {
					case errors.Is(err, data.ErrRecordNotFound):
						count = 0
					default:
						return nil, err
					}
				}

				if status == data.ProductStatusActive && count == 0 {
				}
			}
		default:
			v.AddError("status", "must be a valid status")

		}
	}

	if !v.Valid() {
		return nil, &data.ValidationErrors{
			Fields: v.Errors,
		}
	}

	product, err := s.productRepository.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int32) ([]*data.Product, error) {
	products, err := s.productRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetVariant(ctx context.Context, id int64) (*data.Variant, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	variant, err := s.variantRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *ProductService) CreateVariant(ctx context.Context, productID int64, input *CreateVariantInput) (*data.Variant, error) {
	if productID < 1 {
		return nil, data.ErrRecordNotFound
	}

	v := validator.New()

	validateName(v, input.Name)
	validateSlug(v, input.Slug)
	validateVariantPrice(v, input.Price)

	if !v.Valid() {
		return nil, &data.ValidationErrors{
			Fields: v.Errors,
		}
	}

	product, err := s.productRepository.Get(ctx, productID)
	if err != nil {
		return nil, err
	}

	input.ProductID = product.ID
	variant, err := s.variantRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *ProductService) UpdateVariant(ctx context.Context, id int64, input *UpdateVariantInput) (*data.Variant, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	v := validator.New()

	if input.Name != nil {
		validateName(v, *input.Name)
	}

	if input.Price != nil {
		validateVariantPrice(v, *input.Price)
	}

	if !v.Valid() {
		return nil, &data.ValidationErrors{
			Fields: v.Errors,
		}
	}

	variant, err := s.variantRepository.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *ProductService) DeleteVariant(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.variantRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) ListProductVariants(ctx context.Context, productID int64) ([]*data.Variant, error) {
	if productID < 1 {
		return nil, data.ErrRecordNotFound
	}

	variants, err := s.variantRepository.ListByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return variants, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.productRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func validateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(name) <= 255, "name", "must be at most 255 characters long")
}

func validateSlug(v *validator.Validator, slug string) {
	v.Check(slug != "", "slug", "must be provided")
	v.Check(validator.Matches(slug, validator.HyphenatedRegex), "slug", "must be a valid slug")
	v.Check(utf8.RuneCountInString(slug) <= 128, "slug", "must be at most 128 characters long")
}

func validateStatus(v *validator.Validator, status data.ProductStatus) {
	switch status {
	case data.ProductStatusActive, data.ProductStatusArchived, data.ProductStatusDraft:
	default:
		v.AddError("status", "must be a valid status")
	}
}

func validateVariantPrice(v *validator.Validator, price decimal.Decimal) {
	v.Check(price.GreaterThanOrEqual(decimal.NewFromInt(1)), "price", "must be greater than 0")
}
