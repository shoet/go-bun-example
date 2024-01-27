package repository_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/go-bun-example/entities"
	"github.com/shoet/go-bun-example/infrastracture/repository"
	"github.com/shoet/go-bun-example/testutil"
	"github.com/uptrace/bun"
)

func Test_BookRepository_CreateBook(t *testing.T) {
	type args struct {
		book    *entities.Book
		prepare func(ctx context.Context, tx *bun.Tx) error
	}
	type wants struct {
		book *entities.Book
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
				prepare: nil,
				book: &entities.Book{
					Title: "test01",
					Author: &entities.Author{
						Name: "test01",
					},
				},
			},
			wants: wants{
				book: &entities.Book{
					Title: "test01",
					Author: &entities.Author{
						Name: "test01",
					},
				},
				err: nil,
			},
		},
		{
			name: "success2",
			args: args{
				prepare: func(ctx context.Context, tx *bun.Tx) error {
					author := &entities.Author{
						Name: "test01",
					}
					if _, err := tx.NewInsert().Model(author).Exec(ctx); err != nil {
						return err
					}
					return nil
				},
				book: &entities.Book{
					Title: "test01",
					Author: &entities.Author{
						Name: "test01",
					},
				},
			},
			wants: wants{
				book: &entities.Book{
					Title: "test01",
					Author: &entities.Author{
						Name: "test01",
					},
				},
				err: nil,
			},
		},
	}
	bunDB, closer, err := testutil.ConnectBunDBForTest(t)
	if err != nil {
		t.Fatalf("failed to connect bun db: %v", err)
	}
	t.Cleanup(func() { closer() })
	ctx := context.Background()
	bookRepository, err := repository.NewBookRepository()
	if err != nil {
		t.Fatalf("failed to create book repository: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.DoInTXForTest(t, ctx, bunDB, func(ctx context.Context, tx *bun.Tx) error {
				if tt.args.prepare != nil {
					if err := tt.args.prepare(ctx, tx); err != nil {
						t.Fatalf("failed to prepare: %v", err)
					}
				}
				bookID, err := bookRepository.CreateBook(ctx, tx, tt.args.book)
				if err != nil {
					t.Fatalf("failed to create book: %v", err)
				}

				var got entities.Book
				if err := tx.
					NewSelect().
					Model(&got).
					Where("book.id = ?", bookID).
					Relation("Author").
					Scan(ctx); err != nil {
					t.Fatalf("failed to get book: %v", err)
				}

				cmpoptsBook := cmpopts.IgnoreFields(entities.Book{}, "ID", "AuthorID")
				cmpoptsAuthor := cmpopts.IgnoreFields(entities.Author{}, "ID")
				if diff := cmp.Diff(tt.wants.book, &got, cmpoptsBook, cmpoptsAuthor); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
				return nil
			})
		})
	}

}
