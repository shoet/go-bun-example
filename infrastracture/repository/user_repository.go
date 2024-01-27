package repository

import (
	"context"

	"github.com/shoet/go-bun-example/entities"
	"github.com/uptrace/bun"
)

type UserRepository struct {
}

func NewUserRepository() (*UserRepository, error) {
	return &UserRepository{}, nil
}

func (u *UserRepository) GetUsers(ctx context.Context, tx *bun.Tx, limit int) ([]*entities.User, error) {
	var users []*entities.User
	if err := tx.NewSelect().Model((*entities.User)(nil)).Limit(limit).Scan(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
