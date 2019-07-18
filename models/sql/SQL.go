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
	Data   *User  `json:"data, omitempty"`
	Status Status `json:"status, omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Auth     string `json:"auth"`
	Password string `json:"password"`
	Email    string `json:"email"`
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
		id       int
		login    string
		password string
		auth     string
		email    string
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
		errRow := row.Scan(&id, &login, &password, &auth, &email)
		if errRow != nil {
			log.Fatal(errRow)
		}
		object := User{
			ID:       id,
			Login:    login,
			Password: password,
			Auth:     auth,
			Email:    email,
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
		id       int
		_login   string
		auth     string
		password string
		email    string
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
		errRow := row.Scan(&id, &_login, &password, &auth, &email)
		if errRow != nil {
			log.Fatal(errRow)
		}
		return Response{
			Status: Status{
				Success: true,
				Message: "success",
			},
			Data: &User{
				ID:       id,
				Login:    _login,
				Password: password,
				Auth:     auth,
				Email:    email,
			},
		}
	}
	err = row.Err()
	if err != nil {
		log.Println("HERE")
		log.Fatal(err)
	}
	var e = "no such user"
	return Response{
		Status: Status{
			Success: false,
			Message: e,
		},
	}
}

func (export Export) InsertData(login string, password string, email string) (data []byte) {
	auth := String(30)
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	sqlStatement := fmt.Sprintf(`INSERT  INTO users (first_name, password, auth, email) VALUES ("%s", "%s", "%s", "%s")`, login, hashedPassword, auth, email)
	fmt.Println("############# trying to insert data #############")
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}

	var _, e = db.Exec(sqlStatement)
	if e != nil {
		return toJSON(Response{
			Status: Status{
				Success: false,
				Message: e.Error(),
			},
		})
	}

	return toJSON(getUserData(login))
}

func loginHandler(login string, password string) Response {
	fmt.Println(login)
	user := getUserData(login)
	if user.Status.Success == false {
		return Response{
			Status: Status{
				Success: user.Status.Success,
				Message: user.Status.Message,
			},
		}
	}
	pass := []byte(password)
	hashedPassword := []byte(user.Data.Password)
	err := bcrypt.CompareHashAndPassword(hashedPassword, pass)
	if err != nil {
		return Response{
			Status: Status{
				Success: false,
				Message: "incorrect password",
			},
		}
	}
	//if user.Data.Auth != auth {
	//	return Response{
	//		Status: Status{
	//			Success: false,
	//			Message: "incorrect auth",
	//		},
	//	}
	//}
	return user
}

func updateData(id string, login string, password string, email string, auth string) Response {
	oldUser := getUserData(login)
	if oldUser.Data.Auth != auth {
		return Response{
			Status: Status{
				Success: false,
				Message: "Incorrect auth",
			},
		}
	}
	newAuth := String(30)
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	sqlStatement := `UPDATE users
					 SET first_name = $2, password = $3, email = $4, auth = $5
					 WHERE id = $1`
	db, err := sql.Open("sqlite3", "./data/test.db")
	if err != nil {
		return Response{
			Status: Status{
				Message: err.Error(),
				Success: false,
			},
		}
	}

	var d, e = db.Exec(sqlStatement, login, hashedPassword, email, newAuth, id)
	if e != nil {
		return Response{
			Status: Status{
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
		return getUserData(login)
	}
	return Response{
		Status: Status{
			Message: "error",
			Success: false,
		},
	}
}

func (export Export) GetUserData(login string) (data []byte) {
	var result = getUserData(login)
	return toJSON(result)
}

func (export Export) UpdateUserData(id string, login string, password string, email string, auth string) (data []byte) {
	var result = updateData(id, login, password, email, auth)
	return toJSON(result)
}

func (export Export) Login(login string, password string) []byte {
	return toJSON(loginHandler(login, password))
}

func toJSON(object interface{}) (data []byte) {
	js, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return js
}
