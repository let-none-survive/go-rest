package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-rest/models/sql"
	"net/http"
)

type Export struct {
}

var SQL sql.Export

func userHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query()["login"][0]
	var result = SQL.GetUserData(login)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query()["login"][0]
	password := r.URL.Query()["password"][0]
	var result = SQL.InsertData(login, password)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func (export Export) Serve() {
	router := mux.NewRouter()
	router.HandleFunc("/users", userHandler).Methods("GET")
	router.HandleFunc("/user", insertHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println("Server is listening... http://localhost:8181")
	_ = http.ListenAndServe(":8181", nil)
}
