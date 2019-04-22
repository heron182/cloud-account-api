package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heron182/cloud-account-api/handlers"
	"github.com/heron182/cloud-account-api/schemas"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	schemas.InitDb()

	router := mux.NewRouter()
	router.HandleFunc("/accounts", handlers.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/login", handlers.LoginAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", handlers.GetAccount).Methods("GET")

	fmt.Println("Server running :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
