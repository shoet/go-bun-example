package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("failed to load .env: %s", err))
	}

	var dbConfig DBConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		panic(fmt.Sprintf("failed to process env: %s", err))
	}

	db, dbCloser, err := ConnectMySQL(&dbConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect mysql: %s", err))
	}
	defer dbCloser()

	bunDB := bun.NewDB(db, mysqldialect.New())

	// クエリのログ出力
	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	ctx := context.Background()
	var rnd float64

	// Select a random number.
	if err := bunDB.NewSelect().ColumnExpr("rand()").Scan(ctx, &rnd); err != nil {
		panic(fmt.Sprintf("failed to select rand: %s", err))
	}

	fmt.Println(rnd)
}

func ConnectMySQL(config *DBConfig) (db *sql.DB, closer func() error, err error) {
	mysqlConfig := &mysql.Config{
		Addr:                 fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:                 config.Username,
		Passwd:               config.Password,
		DBName:               config.Database,
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open mysql: %w", err)
	}
	closer = func() error {
		return db.Close()
	}
	if err := db.Ping(); err != nil {
		if err := closer(); err != nil {
			return nil, nil, fmt.Errorf("failed to close mysql: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, closer, nil
}
