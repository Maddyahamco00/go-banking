package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go-banking/internal/middleware"
	"go-banking/internal/pkg/response"
	"go-banking/internal/service"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletSvc *service.WalletService
}

// NewWalletHandler creates a new WalletHandler
func NewWalletHandler(walletSvc *service.WalletService) *WalletHandler {
	return &WalletHandler{walletSvc: walletSvc}
}

// CreateWalletRequest represents the request body for creating a wallet
type CreateWalletRequest struct {
	OwnerID  string `json:"owner_id" binding:"required"`
	Currency string `json:"currency"`
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet for a user
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body CreateWalletRequest true "Wallet creation request"
// @Success 201 {object} response.APIResponse{data=domain.WalletDetails}
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /wallet/create [post]
// @Security BearerAuth
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		response.BadRequest(c, "Invalid owner_id format")
		return
	}

	wallet, err := h.walletSvc.Create(c.Request.Context(), &service.CreateWalletRequest{
		OwnerID:  ownerID,
		Currency: req.Currency,
	})
	if err != nil {
		if err.Error() == "wallet already exists for this owner" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to create wallet: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Wallet created successfully", wallet)
}

// GetWallet godoc
// @Summary Get wallet details
// @Description Get wallet details by owner ID
// @Tags Wallet
// @Produce json
// @Param owner_id path string true "Owner ID"
// @Success 200 {object} response.APIResponse{data=domain.WalletDetails}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /wallet/{owner_id} [get]
// @Security BearerAuth
func (h *WalletHandler) GetWallet(c *gin.Context) {
	ownerIDStr := c.Param("owner_id")
	if ownerIDStr == "" {
		// Try to get from token if not in path
		ownerIDStr, _ = middleware.GetUserID(c)
		if ownerIDStr == "" {
			response.BadRequest(c, "Owner ID required")
			return
		}
	}

	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid owner_id format")
		return
	}

	wallet, err := h.walletSvc.GetWallet(c.Request.Context(), ownerID)
	if err != nil {
		if err.Error() == "wallet not found" {
			response.NotFound(c, "Wallet not found")
			return
		}
		response.InternalError(c, "Failed to get wallet: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Wallet retrieved successfully", wallet)
}

// GetWalletTransactions godoc
// @Summary Get wallet transactions
// @Description Get transaction history for a wallet
// @Tags Wallet
// @Produce json
// @Param owner_id path string true "Owner ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.APIResponse{data=[]domain.Transaction}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /wallet/{owner_id}/transactions [get]
// @Security BearerAuth
func (h *WalletHandler) GetWalletTransactions(c *gin.Context) {
	ownerIDStr := c.Param("owner_id")
	if ownerIDStr == "" {
		ownerIDStr, _ = middleware.GetUserID(c)
		if ownerIDStr == "" {
			response.BadRequest(c, "Owner ID required")
			return
		}
	}

	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid owner_id format")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	txns, total, err := h.walletSvc.GetTransactions(c.Request.Context(), ownerID, limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to get transactions: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", gin.H{
		"transactions": txns,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetWalletLedgerEntries godoc
// @Summary Get wallet ledger entries
// @Description Get ledger entries for a wallet (audit trail)
// @Tags Wallet
// @Produce json
// @Param owner_id path string true "Owner ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /wallet/{owner_id}/ledger [get]
// @Security BearerAuth
func (h *WalletHandler) GetWalletLedgerEntries(c *gin.Context) {
	ownerIDStr := c.Param("owner_id")
	if ownerIDStr == "" {
		ownerIDStr, _ = middleware.GetUserID(c)
		if ownerIDStr == "" {
			response.BadRequest(c, "Owner ID required")
			return
		}
	}

	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid owner_id format")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	entries, total, err := h.walletSvc.GetLedgerEntries(c.Request.Context(), ownerID, limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to get ledger entries: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ledger entries retrieved successfully", gin.H{
		"entries": entries,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// FundWalletRequest represents a wallet funding request
type FundWalletRequest struct {
	Amount        string `json:"amount" binding:"required"`
	Description   string `json:"description"`
	IdempotencyKey string `json:"idempotency_key"`
}

// FundWallet godoc
// @Summary Fund a wallet
// @Description Add funds to a wallet (from external source)
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body FundWalletRequest true "Fund request"
// @Param X-Idempotency-Key header string false "Idempotency key"
// @Success 201 {object} response.APIResponse{data=domain.Transaction}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /wallet/fund [post]
// @Security BearerAuth
func (h *WalletHandler) FundWallet(c *gin.Context) {
	var req FundWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	ownerIDStr, _ := middleware.GetUserID(c)
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// Get wallet to find account ID
	wallet, err := h.walletSvc.GetWallet(c.Request.Context(), ownerID)
	if err != nil {
		response.InternalError(c, "Failed to get wallet: "+err.Error())
		return
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		response.BadRequest(c, "Invalid amount format")
		return
	}

	fundReq := &service.FundRequest{
		OwnerID:        ownerID,
		ToAccountID:     wallet.AccountID,
		Amount:         amount,
		Description:    req.Description,
		IdempotencyKey: req.IdempotencyKey,
	}

	txn, err := h.walletSvc.Fund(c.Request.Context(), fundReq)
	if err != nil {
		response.InternalError(c, "Failed to fund wallet: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Wallet funded successfully", txn)
}