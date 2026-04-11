package api

import (
	"database/sql"
	"net/http"

	db "github.com/Maddyahamco00/go-banking/db"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}

// CreateAccount godoc
// @Summary Create account
// @Description create a new user account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body createAccountRequest true "Account Data"
// @Success 200 {object} db.Account
// @Router /accounts [post]
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetAccount godoc
// @Summary Get account
// @Description get account details by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} db.Account
// @Router /accounts/{id} [get]
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// ListAccounts godoc
// @Summary List accounts
// @Description list all accounts with pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Param page_id query int true "Page ID"
// @Param page_size query int true "Page Size"
// @Success 200 {array} db.Account
// @Router /accounts [get]
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type accountIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type amountRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

// ----------  DEPOSIT ENDPOINT  ----------

// Deposit godoc
// @Summary Deposit funds
// @Description deposit funds into an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param amount body amountRequest true "Deposit Amount"
// @Success 200 {object} db.DepositTxResult
// @Router /accounts/{id}/deposit [post]
func (server *Server) deposit(ctx *gin.Context) {
	var uriReq accountIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bodyReq amountRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := server.store.GetAccount(ctx, uriReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := server.store.DepositTx(ctx, db.DepositTxParams{
		AccountID: uriReq.ID,
		Amount:    bodyReq.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// ----------  WITHDRAW ENDPOINT  ----------

// Withdraw godoc
// @Summary Withdraw funds
// @Description withdraw funds from an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param amount body amountRequest true "Withdraw Amount"
// @Success 200 {object} db.WithdrawTxResult
// @Router /accounts/{id}/withdraw [post]
func (server *Server) withdraw(ctx *gin.Context) {
	var uriReq accountIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bodyReq amountRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := server.store.GetAccount(ctx, uriReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.Balance < bodyReq.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	result, err := server.store.WithdrawTx(ctx, db.WithdrawTxParams{
		AccountID: uriReq.ID,
		Amount:    bodyReq.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// ----------  TRANSACTION HISTORY  ----------

type listEntriesRequest struct {
	AccountID int64 `uri:"id" binding:"required,min=1"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// ListEntries godoc
// @Summary List account entries
// @Description list account transaction entries with pagination
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param page_id query int true "Page ID"
// @Param page_size query int true "Page Size"
// @Success 200 {array} db.Entry
// @Router /accounts/{id}/entries [get]
func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntriesRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entries, err := server.store.ListEntries(ctx, db.ListEntriesParams{
		AccountID: req.AccountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
