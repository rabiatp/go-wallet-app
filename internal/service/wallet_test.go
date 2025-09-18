package service_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/internal/repo"
	"github.com/rabiatp/go-wallet-app/internal/service"
	"github.com/shopspring/decimal"
)

type MockWalletRepo struct {
	balances     map[uuid.UUID]decimal.Decimal
	transactions map[uuid.UUID][]*ent.Transaction
}


// Implement the WalletRepo interface methods with dummy logic
func (m *MockWalletRepo) GetBalance(userID uuid.UUID) (decimal.Decimal, error) {
	// Return a fixed initial balance for testing
	return decimal.NewFromInt(100), nil
}

func (m *MockWalletRepo) UpdateBalance(userID uuid.UUID, amount decimal.Decimal) error {
	// Assume update always succeeds for testing
	return nil
}

func (m *MockWalletRepo) Create(ctx *gin.Context, userID uuid.UUID, balance decimal.Decimal) (*ent.Wallet, error) {
	return &ent.Wallet{ID: uuid.New(), UserID: userID, Balance: balance}, nil
}

func TestDeposit(t *testing.T) {
	walletService := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	userID := uuid.New()
	initialBalance := decimal.NewFromInt(100)

	_, err := walletService.Deposit(ctx, userID, decimal.NewFromInt(-50))
	if err == nil {
		t.Fatal("Deposit with negative amount did not return error")
	}
	
	_, err = walletService.Deposit(ctx, userID, decimal.Zero)
	if err == nil {
		t.Fatal("Deposit with zero amount did not return error")
	}

	depositAmount := decimal.NewFromInt(50)
	newBalance, err := walletService.Deposit(ctx, userID, depositAmount)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}
	expectedBalance := initialBalance.Add(depositAmount)
	if newBalance.Cmp(expectedBalance) != 0 {
		t.Fatalf("Deposit returned incorrect new balance: got %v, want %v", newBalance, expectedBalance)
	}
}
func TestWithdraw(t *testing.T) {
	
	walletService := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	userID := uuid.New()
	initialBalance := decimal.NewFromInt(100)
	
	_, err := walletService.Withdraw(ctx, userID, decimal.NewFromInt(-50))
	if err == nil {
		t.Fatal("Withdraw with negative amount did not return error")
	}
	
	_, err = walletService.Withdraw(ctx, userID, decimal.Zero)		
	if err == nil {
		t.Fatal("Withdraw with zero amount did not return error")
	}			
	
	withdrawAmount := decimal.NewFromInt(50)
	newBalance, err := walletService.Withdraw(ctx, userID, withdrawAmount)		
	if err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}
	expectedBalance := initialBalance.Sub(withdrawAmount)
	if newBalance.Cmp(expectedBalance) != 0 {
		t.Fatalf("Withdraw returned incorrect new balance: got %v, want %v", newBalance, expectedBalance)
	}

	_, err = walletService.Withdraw(ctx, userID, decimal.NewFromInt(200))
	if err == nil {
		t.Fatal("Withdraw with amount greater than balance did not return error")
	}	
}		

func TestCreateWallet(t *testing.T) {
	walletService := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	userID := uuid.New()
	initialBalance := decimal.NewFromInt(100)

	walletID, err := walletService.CreateWallet(ctx, userID, initialBalance)		
	if err != nil {
		t.Fatalf("CreateWallet failed: %v", err)
	}
	if walletID == uuid.Nil {
		t.Fatal("CreateWallet returned nil wallet ID")
	}
}

func TestGetBalance(t *testing.T)	 {
	walletService := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	userID := uuid.New()	
	initialBalance := decimal.NewFromInt(100)

	bal, err := walletService.GetBallance(ctx, userID)
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}	
	if bal.Balance.Cmp(initialBalance) != 0 {
		t.Fatalf("GetBalance returned incorrect balance: got %v, want %v", bal.Balance, initialBalance)
	}
}
func TestListTransactions(t *testing.T) {
	walletService := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	userID := uuid.New()
	
	transactions, err := walletService.ListTransactions(ctx, userID)
	if err != nil {
		t.Fatalf("ListTransactions failed: %v", err)
	}
	if len(transactions) == 0 {
		t.Fatal("ListTransactions returned empty transaction list")
	}		
	
}

func newGinCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}

func TestCreateWallet_And_GetBallance(t *testing.T) {
	svc := service.NewWalletService(&repo.WalletRepo{})
	ctx := newGinCtx()
	uid := uuid.New()

	// create with 100
	wid, err := svc.CreateWallet(ctx, uid, decimal.NewFromInt(100))
	if err != nil {
		t.Fatalf("CreateWallet err: %v", err)
	}
	if wid != uid {
		t.Fatalf("wallet id mismatch: want %s got %s", uid, wid)
	}

	w, err := svc.GetBallance(ctx, uid)
	if err != nil {
		t.Fatalf("GetBallance err: %v", err)
	}
	if !w.Balance.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("balance want 100 got %s", w.Balance.String())
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
    svc := service.NewWalletService(&repo.WalletRepo{})
    ctx := newGinCtx()
    uid := uuid.New()

   
    _, err := svc.CreateWallet(ctx, uid, decimal.NewFromInt(100))
    if err != nil {
        t.Fatalf("CreateWallet err: %v", err)
    }

    // Try to withdraw 150 TL
    _, err = svc.Withdraw(ctx, uid, decimal.NewFromInt(150))
    if err == nil {
        t.Fatal("Withdraw with insufficient funds did not return error")
    }
}

func TestWithdraw_Concurrent(t *testing.T) {
    svc := service.NewWalletService(&repo.WalletRepo{})
    ctx := newGinCtx()
    uid := uuid.New()

    // Create wallet with 100 TL
    _, err := svc.CreateWallet(ctx, uid, decimal.NewFromInt(100))
    if err != nil {
        t.Fatalf("CreateWallet err: %v", err)
    }

    // Simulate concurrent withdrawals
    results := make(chan error, 2)
    go func() {
        _, err := svc.Withdraw(ctx, uid, decimal.NewFromInt(60))
        results <- err
    }()
    go func() {
        _, err := svc.Withdraw(ctx, uid, decimal.NewFromInt(60))
        results <- err
    }()

    var success, fail int
    for i := 0; i < 2; i++ {
        err := <-results
        if err == nil {
            success++
        } else {
            fail++
        }
    }
    if success != 1 || fail != 1 {
        t.Fatalf("Concurrent withdraws: expected 1 success and 1 fail, got %d success, %d fail", success, fail)
    }
}

func TestDepositREST_BalanceGRPC(t *testing.T) {
    svc := service.NewWalletService(&repo.WalletRepo{})
    ctx := newGinCtx()
    uid := uuid.New()

    // Create wallet with 100 TL
    _, err := svc.CreateWallet(ctx, uid, decimal.NewFromInt(100))
    if err != nil {
        t.Fatalf("CreateWallet err: %v", err)
    }

    // REST: Deposit 50 TL
    _, err = svc.Deposit(ctx, uid, decimal.NewFromInt(50))
    if err != nil {
        t.Fatalf("Deposit err: %v", err)
    }

    // gRPC: Get balance
    bal, err := svc.GetBallance(ctx, uid)
    if err != nil {
        t.Fatalf("GetBallance err: %v", err)
    }
    expected := decimal.NewFromInt(150)
    if !bal.Balance.Equal(expected) {
        t.Fatalf("Balance mismatch: want %s got %s", expected.String(), bal.Balance.String())
    }
}