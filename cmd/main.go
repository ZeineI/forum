package main

import (
	"log"
	"os"

	"github.com/ZeineI/forum/internal/database"
	"github.com/ZeineI/forum/internal/server"
	"github.com/ZeineI/forum/internal/session"
)

func main() {

	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	log.SetOutput(file)

	sqlLiteDB := database.SqlLiteDB{}
	if err := sqlLiteDB.Init("forum.db"); err != nil {
		log.Println("DB connection: %v", err)
		return
	}

	sessions := session.InitSession(sqlLiteDB.DB)
	server := server.InitServer(sqlLiteDB, sessions)

	if err := server.Run(); err != nil {
		return
	}
}
