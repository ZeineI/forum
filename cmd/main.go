package main

import (
	"log"
	"os"

	"github.com/ZeineI/forum/internal/database"
	"github.com/ZeineI/forum/internal/server"
	"github.com/ZeineI/forum/internal/session"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	sqlLiteDB := database.SqlLiteDB{}
	sessions := session.InitSession()
	server := server.InitServer(sqlLiteDB, sessions)
	server.Run()
}
