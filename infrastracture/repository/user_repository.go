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

func (u *UserRepository) CreateUser(ctx context.Context, tx *bun.Tx, user *entities.User) error {
	if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUser(
	ctx context.Context, tx *bun.Tx, user *entities.User,
) error {
	if _, err := tx.NewUpdate().Model(user).WherePK().Exec(ctx); err != nil {
		return err
	}
	return nil
}
