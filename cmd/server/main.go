package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rabiatp/go-wallet-app/internal/db"
	httpx "github.com/rabiatp/go-wallet-app/internal/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")

	client := db.NewClient()
	defer client.Close()

	// Hafif bağlantı testi (Ent üzerinden)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// if _, err := client.User.Query().Limit(1).Count(ctx); err != nil {
	// 	// tablo yoksa bile client.Schema.Create zaten oluşturdu; hata gelirse görürüz
	// }

   if _, err := client.User.Query().Limit(1).Count(ctx); err != nil {
    log.Fatalf("db ping: %v", err)
  }

	r := httpx.Router(client);
  srv := &http.Server{ Addr: ":8080", Handler: r }

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
