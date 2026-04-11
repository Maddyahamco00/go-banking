package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Maddyahamco00/go-banking/db"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}

// ---------- LOGIC ----------

// CreateTransfer godoc
// @Summary Create transfer
// @Description transfer funds from one account to another
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body createTransferRequest true "Transfer Data"
// @Success 200 {object} db.TransferTxResult
// @Router /transfers [post]
func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	if fromAccount.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// ---------- VALIDATION ----------

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return account, false
	}

	return account, true
}
