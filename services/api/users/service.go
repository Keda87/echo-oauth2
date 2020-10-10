package users

import (
	"context"
	"github.com/Keda87/echo-oauth2/models"
	"github.com/Keda87/echo-oauth2/services/helper"
	"github.com/jmoiron/sqlx"
)

type userService struct {
	db       *sqlx.DB
	userRepo UserRepositoryInterface
}

func NewService(db *sqlx.DB, userRepo UserRepositoryInterface) *userService {
	return &userService{
		db,
		userRepo,
	}
}

func (svc *userService) Register(ctx context.Context, data *models.UserPayload) (*models.User, error) {
	hashed, err := helper.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	data.Password = hashed
	return svc.userRepo.Insert(ctx, svc.db, data)
}

func (svc *userService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	return svc.userRepo.GetByID(ctx, svc.db, userID)
}
