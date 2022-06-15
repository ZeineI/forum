package main

import (
	"log"
	"os"

	"github.com/ZeineI/forum/internal/database"
	"github.com/ZeineI/forum/internal/server"
	"github.com/ZeineI/forum/internal/session"
)

func init() {
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	log.SetOutput(file)
}

func main() {
	sqlLiteDB := database.SqlLiteDB{}
	sessions := session.InitSession()
	server := server.InitServer(sqlLiteDB, sessions)
	server.Run()
}
