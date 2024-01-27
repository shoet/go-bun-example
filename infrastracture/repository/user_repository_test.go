package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func Test_UpdateRepository_CreateUser(t *testing.T) {
	type args struct {
		updateUser *entities.User
		prepare    func(ctx context.Context, tx *bun.Tx) (*entities.User, error)
	}
	type wants struct {
		user *entities.User
		err  error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "success",
			args: args{
				prepare: func(ctx context.Context, tx *bun.Tx) (*entities.User, error) {
					testUser := &entities.User{Name: "test01"}
					result, err := tx.NewInsert().Model(testUser).Exec(ctx)
					if err != nil {
						return nil, fmt.Errorf("failed to create user: %w", err)
					}
					id, err := result.LastInsertId()
					if err != nil {
						return nil, fmt.Errorf("failed to get last insert id: %w", err)
					}
					return &entities.User{ID: id, Name: "test01"}, nil
				},
				updateUser: &entities.User{Name: "test02"},
			},
			wants: wants{
				user: &entities.User{Name: "test02"},
				err:  nil,
			},
		},
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.DoInTXForTest(t, ctx, bunDB, func(ctx context.Context, tx *bun.Tx) error {
				user, err := tt.args.prepare(ctx, tx)
				if err != nil {
					t.Fatalf("failed to prepare: %v", err)
				}
				updateUser := tt.args.updateUser
				updateUser.ID = user.ID
				if err := userRepository.UpdateUser(ctx, tx, updateUser); err != nil {
					t.Fatalf("failed to update user: %v", err)
				}
				var got entities.User
				if err := tx.
					NewSelect().
					Model((*entities.User)(nil)).
					Where("id = ?", user.ID).
					Scan(ctx, &got); err != nil {
					t.Fatalf("failed to get user: %v", err)
				}

				cmpopts := cmpopts.IgnoreFields(entities.User{}, "ID")
				if diff := cmp.Diff(tt.wants.user, &got, cmpopts); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
				return nil
			})
		})
	}
}
