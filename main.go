package main

import (
	"go-rest/models/server"
	"go-rest/models/sql"
)

var SQL sql.Export
var SERVER server.Export

func main() {
	SQL.ConnectDB()
	SERVER.Serve()
}
