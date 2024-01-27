package repository_test

import (
	"context"
	"testing"

	"github.com/shoet/go-bun-example/entities"
	"github.com/shoet/go-bun-example/infrastracture/repository"
	"github.com/shoet/go-bun-example/testutil"
	"github.com/uptrace/bun"
)

func Test_UserRepository_GetUsers(t *testing.T) {
	bunDB, closer, err := testutil.ConnectBunDBForTest(t)
	if err != nil {
		t.Fatalf("failed to connect bun db: %v", err)
	}
	t.Cleanup(func() { closer() })
	ctx := context.Background()
	userRepository, err := repository.NewUserRepository()
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}
	testutil.DoInTXForTest(t, ctx, bunDB, func(ctx context.Context, tx *bun.Tx) error {
		testUsers := []*entities.User{
			{Name: "test01"}, {Name: "test02"}, {Name: "test03"}, {Name: "test04"}, {Name: "test05"},
			{Name: "test06"}, {Name: "test07"}, {Name: "test08"}, {Name: "test09"}, {Name: "test10"},
			{Name: "test11"}, {Name: "test12"}, {Name: "test13"}, {Name: "test14"}, {Name: "test15"},
		}
		_, err := tx.NewInsert().Model(&testUsers).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create user repository: %v", err)
		}
		users, err := userRepository.GetUsers(ctx, tx, 10)
		if err != nil {
			t.Fatalf("failed to get users: %v", err)
		}

		if len(users) != 10 {
			t.Fatalf("failed to get users: %v", err)
		}

		return nil
	})
}

func Test_UserRepository_CreateUser(t *testing.T) {
	bunDB, closer, err := testutil.ConnectBunDBForTest(t)
	if err != nil {
		t.Fatalf("failed to connect bun db: %v", err)
	}
	t.Cleanup(func() { closer() })
	ctx := context.Background()
	userRepository, err := repository.NewUserRepository()
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}
	testutil.DoInTXForTest(t, ctx, bunDB, func(ctx context.Context, tx *bun.Tx) error {
		user := &entities.User{Name: "test01"}
		if err := userRepository.CreateUser(ctx, tx, user); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		var got entities.User
		if err := tx.
			NewSelect().
			Model((*entities.User)(nil)).
			Where("name = ?", "test01").
			Scan(ctx, &got); err != nil {
			t.Fatalf("failed to get user: %v", err)
		}
		return nil
	})
}
