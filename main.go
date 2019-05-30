package main

import (
	"fmt"
	"learnGO/models/server"
	"learnGO/models/sql"
)

var SQL sql.Export
var SERVER server.Export

func main() {
	fmt.Println("############# connect to data base #############")
	SQL.ConnectDB()
	SERVER.Serve()
}
