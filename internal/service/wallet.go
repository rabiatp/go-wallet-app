package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/internal/repo"
	"github.com/shopspring/decimal"
)

type WalletService struct{ r *repo.WalletRepo }
func NewWalletService(r *repo.WalletRepo) *WalletService {
    return &WalletService{r: r}
}


func (s *WalletService) CreateWallet(ctx *gin.Context, userID uuid.UUID, balance decimal.Decimal) (uuid.UUID, error) {

	w, err := s.r.Create(ctx, userID, balance)	
	if err != nil {
		return w.ID, err
	}
	return w.ID, nil
}

func(s *WalletService)GetBallance(ctx *gin.Context, userID uuid.UUID)(decimal.Decimal, error){
	w, err := s.r.GetBalance(ctx, userID)
	if err != nil {
		return decimal.Zero, err
	}
	return w.Balance, nil
}

func(s *WalletService)Deposit(ctx *gin.Context, userID uuid.UUID, amount decimal.Decimal)(decimal.Decimal, error){
	return s.r.DepositAndLog(ctx, userID, amount)
}

func (s *WalletService) Withdraw(ctx *gin.Context, userID uuid.UUID, amount decimal.Decimal)(decimal.Decimal, error)  {
	return s.r.WithdrawAndLog(ctx, userID, amount)
}

func (s *WalletService) ListTransactions(ctx *gin.Context, userID uuid.UUID) ([]*ent.Transaction, error) {
	return s.r.GetTransactionList(ctx, userID)
}