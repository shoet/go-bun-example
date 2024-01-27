package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoet/go-bun-example/config"
	"github.com/shoet/go-bun-example/infrastracture"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func ConnectBunDBForTest(t *testing.T) (bunDB *bun.DB, closer func() error, err error) {
	t.Helper()
	dbConfig := &config.DBConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "gobun",
		Username: "gobun",
		Password: "gobun",
	}
	db, dbCloser, err := infrastracture.ConnectMySQL(dbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect mysql: %w", err)
	}

	bunDB = bun.NewDB(db, mysqldialect.New())

	// クエリのログ出力
	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return bunDB, dbCloser, nil
}

func DoInTXForTest(
	t *testing.T,
	ctx context.Context,
	bunDB *bun.DB,
	f func(ctx context.Context, tx *bun.Tx) error,
) {
	tx, err := bunDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()
	err = f(ctx, &tx)
	if err != nil {
		t.Fatalf("failed to exec tx: %v", err)
	}
}
