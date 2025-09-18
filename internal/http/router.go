package http

import (
	"net/http"

	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/internal/repo"
	"github.com/rabiatp/go-wallet-app/internal/service"

	"github.com/gin-gonic/gin"
)

func Router(db *ent.Client) *gin.Engine{

	r := gin.Default()

	// health/ping açık
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status":"ok"}) });
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message":"pong"}) });

	//auth
	a := NewAuthServer(db)
	r.POST("/auth/signup", a.SignUp);
	r.POST("/auth/login", a.LogIn);

	// korumalı alan
	wr := repo.NewWalletRepo(db)
	ws := service.NewWalletService(wr)
	wh := NewWalletHandler(ws)

	api := r.Group("/api", AuthMiddleWare())
	{
		api.GET("/wallet/balance", wh.Balance)
		api.POST("/wallet/deposit", wh.Deposit)
		api.POST("/wallet/withdraw", wh.Withdraw)
		api.GET("/wallet/transactions", wh.ListTransactions)
	}

	return r
}
