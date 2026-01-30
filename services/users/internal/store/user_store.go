package store

import (
	"context"
	"time"

	"github.com/amane15/ecom_microservice/internal/dbutils"
	"github.com/amane15/ecom_microservice/services/users/internal/data"
	"github.com/amane15/ecom_microservice/services/users/internal/service"
	"github.com/amane15/ecom_microservice/services/users/internal/store/sqlstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewUserStore(dbpool *pgxpool.Pool) *UserStore {
	return &UserStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (us *UserStore) Get(ctx context.Context, id int64) (*data.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := us.Q.GetUser(ctx, id)
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	user := &data.User{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      string(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return user, nil
}

func (us *UserStore) Create(ctx context.Context, input *service.CreateUserInput) (*data.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := us.Q.CreateUser(ctx, sqlstore.CreateUserParams{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	user := &data.User{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      string(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return user, nil
}

func (us *UserStore) Update(ctx context.Context, id int64, input *service.UpdateUserInput) (*data.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := us.Q.UpdateUser(ctx, sqlstore.UpdateUserParams{
		FirstName: dbutils.PtrToString(input.FirstName),
		LastName:  dbutils.PtrToString(input.LastName),
	})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	user := &data.User{
		ID:        row.ID,
		Email:     row.Email,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      string(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return user, nil
}

func (us *UserStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := us.Q.DeleteUser(ctx, id)
	if err != nil {
		return checkAndHandlePostgresErrors(err)
	}

	return nil
}
