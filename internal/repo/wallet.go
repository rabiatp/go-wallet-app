package repo

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabiatp/go-wallet-app/ent"
	"github.com/rabiatp/go-wallet-app/ent/transaction"
	"github.com/rabiatp/go-wallet-app/ent/wallet"
	"github.com/shopspring/decimal"
)

type WalletRepo struct{ c *ent.Client }

func NewWalletRepo(c *ent.Client) *WalletRepo { return &WalletRepo{c: c} }

func (r *WalletRepo) Create(ctx *gin.Context, userID uuid.UUID, balance decimal.Decimal) (*ent.Wallet, error){
	return r.c.Wallet.Create().SetUserID(userID).SetBalance(balance).SetCreatedAt(time.Now()).Save(ctx)
}

func (r *WalletRepo)  GetBalance(ctx *gin.Context, userID uuid.UUID)(*ent.Wallet, error){
	return r.c.Wallet.Query().Where(wallet.UserID(userID)).Only(ctx)
}

func (r *WalletRepo) DepositAndLog(ctx *gin.Context, userID uuid.UUID, amount decimal.Decimal) (decimal.Decimal, error) {
    tx, err := r.c.Tx(ctx)
    if err != nil {
        return decimal.Zero, err
    }
    
    committed := false
    defer func() {
        if !committed {
            _ = tx.Rollback()
        }
    }()

    w, err := tx.Wallet.
        Query().
        Where(wallet.UserID(userID)).
        ForUpdate().
        Only(ctx)
    if err != nil {
        return decimal.Zero, err
    }

	//Bakiye güncelleniyor
    newBal := w.Balance.Add(amount)
    if _, err = tx.Wallet.
        UpdateOneID(w.ID).
        SetBalance(newBal).
        Save(ctx); err != nil {
        return decimal.Zero, err
    }

    // Transaction log oluştur 
    if _, err = tx.Transaction.
        Create().
        SetWalletID(w.ID).
        SetType(transaction.TypeDEPOSIT).
        SetAmount(amount).
        Save(ctx); err != nil {
        return decimal.Zero, err
    }

    if err = tx.Commit(); err != nil {
        return decimal.Zero, err
    }
    committed = true
    return newBal, nil
}

func (r *WalletRepo) WithdrawAndLog(ctx *gin.Context, userID uuid.UUID, amount decimal.Decimal) (decimal.Decimal, error) {
	tx, err := r.c.Tx(ctx)
    if err != nil {
        return decimal.Zero, err
    }
    
    committed := false
    defer func() {
        if !committed {
            _ = tx.Rollback()
        }
    }()

    w, err := tx.Wallet.
        Query().
        Where(wallet.UserID(userID)).
        ForUpdate().
        Only(ctx)
    if err != nil {
        return decimal.Zero, err
    }

	//Bakiye güncelleniyor
    newBal := w.Balance.Sub(amount)
	if newBal.IsNegative() {
		return decimal.Zero, errors.New("an amount larger than the balance cannot be withdrawn")
	}

	// Transaction log oluştur 
    if _, err = tx.Transaction.
        Create().
        SetWalletID(w.ID).
        SetType(transaction.TypeWITHDRAW).
        SetAmount(amount).
        Save(ctx); err != nil {
        return decimal.Zero, err
    }

    if err = tx.Commit(); err != nil {
        return decimal.Zero, err
    }
    committed = true
	
	return newBal, err
}

func (r *WalletRepo) GetTransactionList(ctx *gin.Context, userID uuid.UUID)([]*ent.Transaction, error) {
    return r.c.Transaction.Query().
    Where(transaction.HasWalletWith(wallet.UserID(userID))). //önce walleti bulup sonra Id syle tek sorguda gözüyor
    Order(ent.Desc(transaction.FieldCreatedAt)).
    All(ctx)
}