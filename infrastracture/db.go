package infrastracture

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/shoet/go-bun-example/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func ConnectMySQL(config *config.DBConfig) (db *sql.DB, closer func() error, err error) {
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

func ConnectBunDB(config *config.DBConfig) (bunDB *bun.DB, closer func() error, err error) {
	db, dbCloser, err := ConnectMySQL(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect mysql: %w", err)
	}

	bunDB = bun.NewDB(db, mysqldialect.New())

	// クエリのログ出力
	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return bunDB, dbCloser, nil
}
