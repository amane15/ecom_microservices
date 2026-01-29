package service

import (
	"context"
	"time"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
)

type CreateCategoryInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type CategoryRepository interface {
	Get(ctx context.Context, id int64) (*data.Category, error)
	Create(ctx context.Context, input *CreateCategoryInput) (*data.Category, error)
	Update(ctx context.Context, id int64, input *UpdateCategoryInput) (*data.Category, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int32) ([]*data.Category, error)
}

type CategoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(r CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: r,
	}
}

func (s CategoryService) GetCategory(ctx context.Context, id int64) (*data.Category, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	category, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s CategoryService) CreateCategory(ctx context.Context, input *CreateCategoryInput) (*data.Category, error) {
	v := validator.New()

	validateName(v, input.Name)
	validateSlug(v, input.Slug)

	if !v.Valid() {
		return nil, &data.ValidationErrors{
			Fields: v.Errors,
		}
	}

	category, err := s.categoryRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s CategoryService) UpdateCategory(ctx context.Context, id int64, input *UpdateCategoryInput) (*data.Category, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	v := validator.New()

	if input.Name != nil {
		validateName(v, *input.Name)
	}

	category, err := s.categoryRepository.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.categoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s CategoryService) ListCategories(ctx context.Context, limit, offset int32) ([]*data.Category, error) {
	categories, err := s.categoryRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
