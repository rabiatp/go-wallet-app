package http

import (
	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/internal/service"

	"github.com/gin-gonic/gin"
)

func Router(db *ent.Client, ws *service.WalletService) *gin.Engine {
	r := gin.Default()

	// auth server, ws’ye ihtiyaç duyuyor (signup’ta wallet açmak için)
	auth := NewAuthServer(db, ws)
	r.POST("/auth/signup", auth.SignUp)
	r.POST("/auth/login", auth.LogIn)

	// wallet handler
	wh := NewWalletHandler(ws)

	api := r.Group("/api", AuthMiddleWare()) // JWT middleware uid set etmeli
	{
		api.GET("/wallet/balance", wh.Balance)
		api.POST("/wallet/deposit", wh.Deposit)
		api.POST("/wallet/withdraw", wh.Withdraw)
		api.GET("/wallet/transactions", wh.ListTransactions) // varsa
	}

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	return r
}
