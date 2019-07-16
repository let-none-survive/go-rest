package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"time"
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
	Auth      string `json:"auth"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func (export Export) ConnectDB() {
	if _, err := os.Stat("./data/test.db"); os.IsNotExist(err) {
		fmt.Println("create folder")
		_ = os.MkdirAll(`./data`, 0755)
		_, _ = os.Create("./data/test.db")
	}
	createTable()
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
func createTable() {
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("############# starting create table #############")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( id INTEGER NOT NULL PRIMARY KEY, `first_name` VARCHAR(255) NOT NULL UNIQUE, `password` VARCHAR(255) NOT NULL, `auth` VARCHAR(255) NOT NULL, `email` VARCHAR(255) UNIQUE NOT NULL)")
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
		email     string
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
		errRow := row.Scan(&id, &firstName, &password, &auth, &email)
		if errRow != nil {
			log.Fatal(errRow)
		}
		object := User{
			ID:        id,
			FirstName: firstName,
			Password:  password,
			Auth:      auth,
			Email:     email,
		}
		fmt.Println(object)
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
		auth      string
		password  string
		email     string
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
		errRow := row.Scan(&id, &firstName, &password, &auth, &email)
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
				Email:     email,
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

func (export Export) InsertData(firstName string, password string, email string) (data []byte) {
	auth := String(30)
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	sqlStatement := fmt.Sprintf(`INSERT  INTO users (first_name, password, auth, email) VALUES ("%s", "%s", "%s", "%s")`, firstName, hashedPassword, auth, email)
	fmt.Println("############# trying to insert data #############")
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}

	var _, e = db.Exec(sqlStatement)
	if e != nil {
		return toJSON(Response{
			Status: &Status{
				Success: false,
				Message: e.Error(),
			},
		})
	}

	return toJSON(getUserData(firstName))
}

func updateData(id string, firstName string, password string, email string) Response {
	auth := String(30)
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	sqlStatement := `UPDATE users
					 SET first_name = $2, password = $3, email = $4, auth = $5
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

	var d, e = db.Exec(sqlStatement, firstName, hashedPassword, email, auth, id)
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

func (export Export) UpdateUserData(id string, firstName string, password string, email string) (data []byte) {
	var result = updateData(id, firstName, password, email)
	return toJSON(result)
}

func toJSON(object interface{}) (data []byte) {
	js, err := json.Marshal(&object)
	if err != nil {
		panic(err)
	}
	return js
}
