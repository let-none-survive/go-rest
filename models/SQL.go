package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Export struct {
}

type User struct {
	id        int
	firstName string
	password  string
}

type Response struct {
	Message string `json:"message"`
	Data    Data   `json:"data, omitempty"`
}

type Data struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"password"`
}

func (export Export) ConnectDB() {
	if _, err := os.Stat("./data/test.db"); os.IsNotExist(err) {
		fmt.Println("create folder")
		_ = os.MkdirAll(`./data`, 0755)
		_, _ = os.Create("./data/test.db")
	}
	createTable()
}

func createTable() {
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("############# starting create table #############")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( id INTEGER NOT NULL PRIMARY KEY, `first_name` VARCHAR(255) NOT NULL UNIQUE, `password` VARCHAR(255) NOT NULL)")
	if err != nil {
		fmt.Println("############# fatal error #############")
		log.Fatal(err)
	}
	fmt.Println("############# database connected #############")
}

func getUserDate(login string) Response {
	var (
		id        int
		firstName string
		password  string
	)
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	sqlStatement := "SELECT * from users  WHERE first_name = ($1) LIMIT 1;"
	row, err := db.Query(sqlStatement, login)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		errRow := row.Scan(&id, &firstName, &password)
		if errRow != nil {
			log.Fatal(errRow)
		}

		return Response{
			Message: "success",
			Data: Data{
				ID:        id,
				FirstName: firstName,
				LastName:  password,
			},
		}
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}
	return Response{}
}

func (export Export) InsertData(firstName string, password string) (data []byte) {
	sqlStatement := fmt.Sprintf(`
INSERT  INTO users (first_name, password)
VALUES ("%s", "%s")`, firstName, password)
	fmt.Println("############# trying to insert data #############")
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	var _, e = db.Exec(sqlStatement)
	if e != nil {
		fmt.Println(e)
		return toJSON(Response{Message: "error"})
	}

	var newUser = getUserDate(firstName)
	return toJSON(newUser)
}

func (export Export) GetUserData(login string) (data []byte) {
	var result = getUserDate(login)
	return toJSON(result)
}

func toJSON(object interface{}) (data []byte) {
	js, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return js
}
