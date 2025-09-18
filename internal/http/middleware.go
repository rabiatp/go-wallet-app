package http

import (
	"github.com/rabiatp/go-wallet-app/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc{

	return func(ctx *gin.Context) {
		h := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"});
			return 
		}
		token := strings.TrimSpace(h[7:])
		uid, err := auth.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"});
			return 
		}

		ctx.Set("uid", uid);
		ctx.Next();
	}
}
