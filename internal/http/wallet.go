package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabiatp/go-wallet-app/internal/service"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletHandler struct{ s service.WalletService } 
type walletDTO struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
}
func NewWalletHandler(s *service.WalletService) *WalletHandler { return &WalletHandler{s: *s} }

func (h *WalletHandler) Deposit(ctx *gin.Context) {
	var in walletDTO

	if err := ctx.ShouldBindJSON(&in); err != nil {
	
		status.Errorf(codes.Unavailable, err.Error())
		return
	}
	
	uidStr, ok := ctx.Get("uid")
	if !ok {
		status.Errorf(codes.Unauthenticated, "missing uid in context")
		return
	}
	userID, err := uuid.Parse(uidStr.(string))
	if err != nil {
		status.Errorf(codes.InvalidArgument, "invalid uid")
		return
	}

	newBal, err := h.s.Deposit(ctx, userID, in.Amount)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": newBal})
}

func (h *WalletHandler) Withdraw(ctx *gin.Context) {
	var in walletDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		status.Errorf(codes.Unavailable, err.Error())
		return
	}

	uidStr, ok := ctx.Get("uid")
	if !ok {
		status.Errorf(codes.Unauthenticated, "missing uid in context")
		return
	}
	userID, err := uuid.Parse(uidStr.(string))
	if err != nil {
		status.Errorf(codes.InvalidArgument, "invalid uid")
		return
	}
	newBal, err := h.s.Withdraw(ctx, userID, in.Amount)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": newBal})
}

func (h *WalletHandler) ListTransactions(ctx *gin.Context) {
	uidStr := ctx.GetString("uid")

    userID, _ := uuid.Parse(uidStr) 

	list, err := h.s.ListTransactions(ctx, userID)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"transaction": list})
}

func (h *WalletHandler) Balance(ctx *gin.Context) {
	uidStr := ctx.GetString("uid")

    userID, _ := uuid.Parse(uidStr) 

	bal, err := h.s.GetBallance(ctx, userID)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"balance": bal.Balance})

}


