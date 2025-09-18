package grpcserver

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rabiatp/go-wallet-app/internal/service"
	walletv1 "github.com/rabiatp/go-wallet-app/proto/gen/wallet/v1"
)

type WalletServer struct {
	walletv1.UnimplementedWalletServer
	Svc *service.WalletService
}

func NewWalletServer(svc *service.WalletService) *WalletServer {
	return &WalletServer{Svc: svc}
}

func (s *WalletServer) GetBalance(ctx *gin.Context, req *walletv1.GetBalanceRequest) (*walletv1.GetBalanceResponse, error) {
	uid, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	bal, err := s.Svc.GetBallance(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &walletv1.GetBalanceResponse{Balance: bal.Balance.String()}, nil
}

func (s *WalletServer) Deposit(ctx *gin.Context, req *walletv1.DepositRequest) (*walletv1.DepositResponse, error) {
	uid, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}
	newBal, err := s.Svc.Deposit(ctx, uid, amount)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &walletv1.DepositResponse{NewBalance: newBal.String()}, nil
}

func (s *WalletServer) Withdraw(ctx *gin.Context, req *walletv1.WithdrawRequest) (*walletv1.WithdrawResponse, error) {
	uid, err := uuid.Parse(req.GetUserId())	
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}
	newBal, err := s.Svc.Withdraw(ctx, uid, amount)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &walletv1.WithdrawResponse{NewBalance: newBal.String()}, nil
}	
func (s *WalletServer) ListTransactions(ctx *gin.Context, req *walletv1.ListTransactionsRequest) (*walletv1.ListTransactionsResponse, error) {
	uid, err := uuid.Parse(req.GetUserId())		
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	list, err := s.Svc.ListTransactions(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	resp := &walletv1.ListTransactionsResponse{Transactions: make([]*walletv1.Transaction, 0, len(list))}
	for _, t := range list {
		resp.Transactions = append(resp.Transactions, &walletv1.Transaction{
			Id:        t.ID.String(),
			Type:      t.Type.String(),      // "deposit"/"withdraw"
			Amount:    t.Amount.String(),    // decimal -> string
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}