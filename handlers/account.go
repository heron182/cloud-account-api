package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/heron182/cloud-account-api/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateAccount handles creation of a new Account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account schemas.Account
	json.NewDecoder(r.Body).Decode(&account)

	if err := account.Create(); err != nil {
		log.Fatal(err)
	} else {

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(&account)

	}

}

// GetAccount gets an account by ID
func GetAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var account schemas.Account

	if err := account.Get(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	w.Header().Add("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

// LoginAccount check if an account with email/password exists
func LoginAccount(w http.ResponseWriter, r *http.Request) {
	var account schemas.Account
	json.NewDecoder(r.Body).Decode(&account)

	err := account.CheckCredentials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": account.Email,
		"iat":   time.Now().Unix(),
	})

	signedToken, _ := token.SignedString([]byte(os.Getenv("MYSECRET")))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "token": "` + signedToken + `" }`))
}
