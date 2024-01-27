package repository_test

import (
	"context"
	"testing"

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
	testutil.DoInTXForTest(t, ctx, bunDB, func(ctx context.Context, tx *bun.Tx) error {
		prepareQuery := `
		INSERT INTO users (name)
		VALUES 
			("test01"), ("test02"), ("test03"), ("test04"), ("test05"), 
			("test06"), ("test07"), ("test08"), ("test09"), ("test10"),
			("test11"), ("test12"), ("test13"), ("test14"), ("test15"), 
			("test16"), ("test17"), ("test18"), ("test19"), ("test20")
		;
		`
		if _, err := tx.Exec(prepareQuery); err != nil {
			t.Fatalf("failed to prepare query: %v", err)
		}

		userRepository, err := repository.NewUserRepository()
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
