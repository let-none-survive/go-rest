package main

import (
	"fmt"
	"learnGO/models"
	"learnGO/server"
)

var SQL models.Export
var SERVER server.Export

func main() {
	fmt.Println("############# connect to data base #############")
	SQL.ConnectDB()
	SERVER.Serve()
}
