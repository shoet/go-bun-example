package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/shoet/go-bun-example/config"
	"github.com/shoet/go-bun-example/infrastracture"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("failed to load .env: %s", err))
	}

	dbConfig, err := config.NewDBConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load db config: %s", err))
	}

	bunDB, dbCloser, err := infrastracture.ConnectBunDB(dbConfig)
	defer dbCloser()

	ctx := context.Background()
	var rnd float64

	// Select a random number.
	if err := bunDB.NewSelect().ColumnExpr("rand()").Scan(ctx, &rnd); err != nil {
		panic(fmt.Sprintf("failed to select rand: %s", err))
	}

	fmt.Println(rnd)
}
