package users

import (
	"context"
	"github.com/Keda87/echo-oauth2/models"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	GetTableName() string
	Insert(ctx context.Context, db *sqlx.DB, data *models.UserPayload) (*models.User, error)
	GetByID(ctx context.Context, db *sqlx.DB, userID uint) (*models.User, error)
}

type UserServiceInterface interface {
	Register(ctx context.Context, data *models.UserPayload) (*models.User, error)
	GetByID(ctx context.Context, userID uint) (*models.User, error)
}
