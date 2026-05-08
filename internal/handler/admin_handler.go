package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-banking/internal/repository"
	"go-banking/internal/pkg/response"
)

// AdminHandler handles admin API requests
type AdminHandler struct {
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
	walletRepo      *repository.WalletRepository
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
	walletRepo *repository.WalletRepository,
) *AdminHandler {
	return &AdminHandler{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
	}
}

// GetTransactions godoc
// @Summary List all transactions (admin)
// @Description Get all transactions with pagination
// @Tags Admin
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /admin/transactions [get]
// @Security BearerAuth
func (h *AdminHandler) GetTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	txns, err := h.transactionRepo.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to get transactions: "+err.Error())
		return
	}

	total, _ := h.transactionRepo.Count(c.Request.Context())

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", gin.H{
		"transactions": txns,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetWallets godoc
// @Summary List all wallets (admin)
// @Description Get all wallets with pagination
// @Tags Admin
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /admin/wallets [get]
// @Security BearerAuth
func (h *AdminHandler) GetWallets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	accounts, err := h.accountRepo.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to get wallets: "+err.Error())
		return
	}

	total, _ := h.accountRepo.Count(c.Request.Context())

	response.Success(c, http.StatusOK, "Wallets retrieved successfully", gin.H{
		"wallets": accounts,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetLoans godoc
// @Summary List all loans (admin)
// @Description Get all loans with pagination
// @Tags Admin
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /admin/loans [get]
// @Security BearerAuth
func (h *AdminHandler) GetLoans(c *gin.Context) {
	// TODO: Implement when loan service is ready
	response.Success(c, http.StatusOK, "Loans retrieved successfully", gin.H{
		"loans":  []interface{}{},
		"total":  0,
		"limit":  20,
		"offset": 0,
	})
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags Health
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /health [get]
func (h *AdminHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"status":  "healthy",
		"service": "go-banking",
		"version": "1.0.0",
	})
}