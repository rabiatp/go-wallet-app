package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/rabiatp/go-wallet-app/internal/db"
	httpx "github.com/rabiatp/go-wallet-app/internal/http"
	"github.com/rabiatp/go-wallet-app/internal/repo"
	"github.com/rabiatp/go-wallet-app/internal/service"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load("../.env")

	client := db.NewClient()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.User.Query().Limit(1).Count(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	//service
	wr := repo.NewWalletRepo(client)
	ws := service.NewWalletService(wr)

	//HTTP router
	r := httpx.Router(client, ws) 

	// gRPC server
	gs := grpc.NewServer() 
	//walletv1.RegisterWalletServer(gs, grpcserver.NewWalletServer(ws)) 

	//gRPC: :9090
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("gRPC listen: %v", err)
		}
		log.Println("[gRPC] listening on :9090")
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("gRPC serve: %v", err)
		}
	}()

	// HTTP: :8080
	srv := &http.Server{Addr: ":8080", Handler: r}
	log.Println("[HTTP] listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
