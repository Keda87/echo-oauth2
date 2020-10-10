package users

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Keda87/echo-oauth2/models"
)

type userRepository struct {
}

func NewRepository() *userRepository {
	return &userRepository{}
}
func (repo *userRepository) GetTableName() string {
	return "users"
}

func (repo *userRepository) Insert(ctx context.Context, db *sqlx.DB, data *models.UserPayload) (*models.User, error) {
	query, args, _ := sq.Insert(repo.GetTableName()).
		Columns("email", "password").
		Values(data.Email, data.Password).
		Suffix("RETURNING id").
		ToSql()
	query = db.Rebind(query)

	var insertedID uint
	err := db.GetContext(ctx, &insertedID, query, args...)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(ctx, db, insertedID)
}

func (repo *userRepository) GetByID(ctx context.Context, db *sqlx.DB, userID uint) (*models.User, error) {
	var result models.User

	query, args, _ := sq.Select("id", "email", "password").
		From(repo.GetTableName()).
		Where("id = ?", userID).
		ToSql()
	query = db.Rebind(query)

	err := db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
