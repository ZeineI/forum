package main

func main() {
	sqlLiteDB := database.SqlLiteDB{}
	sessions := session.InitSession()
	server := server.InitServer(sqlLiteDB, sessions)
	server.Run()
}
