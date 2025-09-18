package db

import (
	"context"
	"log"
	"os"

	"github.com/rabiatp/go-wallet-app/ent"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq" // Postgres driver
)

func NewClient() *ent.Client {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// .env yüklenmezse bile çalışsın
		dsn = "postgres://postgres:254798@localhost:5433/wallet_db?sslmode=disable"
	}

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	if err := client.Schema.Create( 
		context.Background(),
    	schema.WithDropColumn(true),
    	schema.WithDropIndex(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
