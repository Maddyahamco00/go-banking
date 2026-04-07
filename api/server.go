package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/Maddyahamco00/go-banking/db"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/accounts/:id/deposit", server.deposit)
	router.POST("/accounts/:id/withdraw", server.withdraw)

	router.POST("/transfers", server.createTransfer)

	router.GET("/accounts/:id/entries", server.listEntries)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}