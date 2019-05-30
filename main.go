package main

import (
	"fmt"
	"go-rest/models/server"
	"go-rest/models/sql"
)

var SQL sql.Export
var SERVER server.Export

func main() {
	fmt.Println("############# connect to data base #############")
	SQL.ConnectDB()
	SERVER.Serve()
}
