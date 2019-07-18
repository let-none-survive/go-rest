package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-rest/models/sql"
	"net/http"
)

type Export struct {
}

type User struct {
	Password string `json:"password"`
	Login    string `json:"login"`
	Auth     string `json:"auth"`
}

type FullUser struct {
	User
	Email string `json:"email"`
}

var SQL sql.Export

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	if user == "all" {
		var result = SQL.GetAllUsersData()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(result)
		return
	}
	var result = SQL.GetUserData(user)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)

	return

}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FullUser
	err := decoder.Decode(&user)

	if err != nil {
		fmt.Println("err")
		panic(err)
	}
	fmt.Println(user)
	var result = SQL.InsertData(user.Login, user.Password, user.Email)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)

	var user FullUser
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}
	result := SQL.UpdateUserData(id, user.Login, user.Password, user.Email, user.Auth)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(SQL.Login(user.Login, user.Password))
}

func (export Export) Serve() {
	router := mux.NewRouter()
	router.HandleFunc("/users/{user}", userHandler).Methods("GET")
	router.HandleFunc("/user/{id}", updateHandler).Methods("PATCH")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/user", insertHandler).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Server is listening... http://localhost:8181")
	fmt.Println("Routes: ")
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	_ = http.ListenAndServe(":8181", nil)
}
