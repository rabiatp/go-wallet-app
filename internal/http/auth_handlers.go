package http

import (
	"net/http"
	"time"

	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/ent/user"
	"github.com/rabiatp/go-wallet-app/internal/auth"
	"github.com/rabiatp/go-wallet-app/internal/repo"
	"github.com/rabiatp/go-wallet-app/internal/service"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

type AuthServer struct{ db *ent.Client }
type WalletService struct{ s *service.WalletService }
func NewAuthServer(db *ent.Client, wsvc *service.WalletService) *AuthServer {
	return &AuthServer{db: db}
}

type signupDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthServer)SignUp (ctx *gin.Context)  {
	
	var in signupDTO

	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); 
		return
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		 ctx.JSON(500, gin.H{"error":"hash"});
		 return
	}

	u, err := a.db.User.Create().SetName(in.Name).SetEmail(in.Email).SetPasswordHash(hash).Save(ctx)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()});
		return
	}
	token, err := auth.GenerateToken(u.ID.String(), 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"});
		return
	}

	//hesap açılımında balance 0 olarak cüzdan oluşturulur
	wRepo := repo.NewWalletRepo(a.db)
	wSvc  := service.NewWalletService(wRepo)  //burası kızarıyor

	if _, err := wSvc.CreateWallet(ctx, u.ID, decimal.Zero); err != nil {
    	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "wallet_create_failed"})
    return
}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user_id": u.ID})
}

func (a *AuthServer)LogIn(ctx *gin.Context)  {
	var in loginDTO
	
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	
	u, err := a.db.User.Query().Where(user.Email(in.Email)).Only(ctx)
	if err != nil { 
		ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid creds"}); 
		return 
	}
	
	if err := auth.CheckPassword(u.PasswordHash, in.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid creds"}); 
		return
	}
  
	token, err := auth.GenerateToken(u.ID.String(), 24*time.Hour)
  	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"});
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user_id": u.ID})
}


