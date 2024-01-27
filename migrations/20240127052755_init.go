package migrations

import (
	"context"
	"fmt"

	"github.com/shoet/go-bun-example/entities"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		if _, err := db.NewCreateTable().Model((*entities.User)(nil)).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table users: %w", err)
		}
		if _, err := db.NewCreateTable().Model((*entities.Book)(nil)).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table books: %w", err)
		}
		if _, err := db.NewCreateTable().Model((*entities.Author)(nil)).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table authors: %w", err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		if _, err := db.NewDropTable().Model((*entities.User)(nil)).IfExists().Exec(ctx); err != nil {
			return fmt.Errorf("failed to drop table users: %w", err)
		}
		if _, err := db.NewDropTable().Model((*entities.Book)(nil)).IfExists().Exec(ctx); err != nil {
			return fmt.Errorf("failed to drop table books: %w", err)
		}
		if _, err := db.NewDropTable().Model((*entities.Author)(nil)).IfExists().Exec(ctx); err != nil {
			return fmt.Errorf("failed to drop table authors: %w", err)
		}
		return nil
	})
}
