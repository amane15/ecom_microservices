package service

import (
	"context"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/users/internal/data"
)

type CreateUserInput struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type UserRepository interface {
	Get(ctx context.Context, id int64) (*data.User, error)
	Create(ctx context.Context, input *CreateUserInput) (*data.User, error)
	Update(ctx context.Context, id int64, input *UpdateUserInput) (*data.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{
		userRepository: r,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*data.User, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) CreateUser(ctx context.Context, input *CreateUserInput) (*data.User, error) {
	v := validator.New()

	validateName(v, input.FirstName, "first_name")
	validateName(v, input.LastName, "last_name")
	validateEmail(v, input.Email)

	if !v.Valid() {
		return nil, &data.ValidationError{
			Fields: v.Errors,
		}
	}

	user, err := s.userRepository.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, input *UpdateUserInput) (*data.User, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	v := validator.New()

	if input.FirstName != nil {
		validateName(v, *input.FirstName, "first_name")
	}
	if input.LastName != nil {
		validateName(v, *input.LastName, "last_name")
	}

	user, err := s.userRepository.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func validateName(v *validator.Validator, name string, key string) {
	v.Check(name != "", key, "must be provided")
	v.Check(utf8.RuneCountInString(name) <= 255, key, "must not be more thann 255 bytes long")
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
