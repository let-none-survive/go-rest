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
	login := r.URL.Query()["login"][0]
	password := r.URL.Query()["password"][0]
	var result = SQL.InsertData(login, password)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	login := r.URL.Query()["login"][0]
	password := r.URL.Query()["password"][0]
	result := SQL.UpdateUserData(id, login, password)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func (export Export) Serve() {
	router := mux.NewRouter()
	router.HandleFunc("/users/{user}", userHandler).Methods("GET")
	router.HandleFunc("/user/{id}", updateHandler).Methods("PATCH")
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
