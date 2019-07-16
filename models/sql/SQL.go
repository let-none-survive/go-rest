package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Gurpartap/async"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"os"
)

type Export struct{}

type Status struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Data   *User   `json:"data, omitempty"`
	Status *Status `json:"status, omitempty"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
	Auth      string `json:"auth"`
}

func (export Export) ConnectDB() {
	if _, err := os.Stat("./data/test.db"); os.IsNotExist(err) {
		fmt.Println("create folder")
		_ = os.MkdirAll(`./data`, 0755)
		_, _ = os.Create("./data/test.db")
	}
	createTable()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func createTable() {
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("############# starting create table #############")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( id INTEGER NOT NULL PRIMARY KEY, `first_name` VARCHAR(255) NOT NULL UNIQUE, `password` VARCHAR(255) NOT NULL, `auth` VARCHAR(255) NOT NULL)")
	if err != nil {
		fmt.Println("############# fatal error #############")
		log.Fatal(err)
	}
	fmt.Println("############# database connected #############")
}

func getAllUsers() []User {
	var (
		id        int
		firstName string
		password  string
		auth      string
	)

	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}

	sqlStatement := "SELECT * from users"
	row, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal(err)
	}

	defer row.Close()

	var result []User

	for row.Next() {
		errRow := row.Scan(&id, &firstName, &password, &auth)
		if errRow != nil {
			log.Fatal(errRow)
		}
		object := User{
			ID:        id,
			FirstName: firstName,
			Password:  password,
			Auth:      auth,
		}
		result = append(result, object)
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (export Export) GetAllUsersData() (data []byte) {
	var result = getAllUsers()
	return toJSON(result)
}

func getUserData(login string) Response {
	var (
		id        int
		firstName string
		password  string
		auth      string
	)
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	sqlStatement := "SELECT * from users  WHERE first_name = $1;"
	row, err := db.Query(sqlStatement, login)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		errRow := row.Scan(&id, &firstName, &password, &auth)
		if errRow != nil {
			log.Fatal(errRow)
		}
		return Response{
			Status: &Status{
				Success: true,
				Message: "success",
			},
			Data: &User{
				ID:        id,
				FirstName: firstName,
				Password:  password,
				Auth:      auth,
			},
		}
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}
	var e = "no such user"
	return Response{
		Status: &Status{
			Success: false,
			Message: e,
		},
	}
}

func (export Export) InsertData(firstName string, password string) (data []byte) {
	auth := randStringRunes(30)
	sqlStatement := fmt.Sprintf(`INSERT  INTO users (first_name, password, auth) VALUES ("%s", "%s", "%s")`, firstName, password, auth)
	fmt.Println("############# trying to insert data #############")
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}

	var _, e = db.Exec(sqlStatement)
	if e != nil {
		fmt.Println(e)
		var errorString = fmt.Sprintf(`user with login %s already exits`, firstName)
		return toJSON(Response{
			Status: &Status{
				Success: false,
				Message: errorString,
			},
		})
	}

	future := async.Any(func() (interface{}, error) {
		return getUserData(firstName), nil
	})

	v, err := future.Await()
	newUser := v.(*User)
	return toJSON(newUser)
}

func updateData(id string, firstName string, password string) Response {
	sqlStatement := `UPDATE users
					 SET first_name = $2, password = $3 
					 WHERE id = $1`
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		return Response{
			Status: &Status{
				Message: err.Error(),
				Success: false,
			},
		}
	}

	var d, e = db.Exec(sqlStatement, firstName, password, id)
	if e != nil {
		return Response{
			Status: &Status{
				Message: e.Error(),
				Success: false,
			},
		}
	}
	count, err := d.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count > 0 {
		return getUserData(firstName)
	}
	return Response{
		Status: &Status{
			Message: "error",
			Success: false,
		},
	}
}

func (export Export) GetUserData(login string) (data []byte) {
	var result = getUserData(login)
	return toJSON(result)
}

func (export Export) UpdateUserData(id string, firstName string, password string) (data []byte) {
	var result = updateData(id, firstName, password)
	return toJSON(result)
}

func toJSON(object interface{}) (data []byte) {
	js, err := json.Marshal(&object)
	if err != nil {
		panic(err)
	}
	return js
}
