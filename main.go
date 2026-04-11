// @title Banking API
// @version 1.0
// @description This is a fintech API system
// @host localhost:8080
// @BasePath /
package main

import (
	"database/sql"
	"log"

	"github.com/Maddyahamco00/go-banking/api"
	db "github.com/Maddyahamco00/go-banking/db"
	"github.com/Maddyahamco00/go-banking/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.ServerAddr); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
