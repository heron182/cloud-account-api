package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"github.com/heron182/cloud-account-api/schemas"
	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("-- No .env file encountered in the root project folder")
	}

	schemas.InitDb(os.Getenv("DATABASE_URI"))
}

func tearDown() {
	schemas.Db.Database("accounts").Drop(context.Background())
}

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())

	tearDown()
}

func makeRequest(req *http.Request, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)

	return rr
}

func TestCreateAccount(t *testing.T) {
	defer tearDown()
	payload := []byte(`{"name": "John", "email": "john@createme.com", "password": "123pass"}`)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(payload))
	rr := makeRequest(req, CreateAccount)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("TestCreateAccount failed. Expected status %d, got %d", http.StatusCreated, status)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	id, _ := primitive.ObjectIDFromHex(response["id"].(string))

	var acc schemas.Account
	if err := acc.Get(id); err != nil {
		t.Errorf("TestCreateAccount failed - %s", err)
	}

}

func TestLoginAccount(t *testing.T) {
	defer tearDown()
	plainPassword := "mYpass"

	acc := schemas.Account{Name: "Johnny", Password: plainPassword, Email: "jonny@email.com"}
	acc.Create()

	payload := []byte(fmt.Sprintf(`{"email": "jonny@email.com", "password": "%s"}`, plainPassword))
	req, _ := http.NewRequest("POST", "/accounts/login", bytes.NewReader(payload))
	rr := makeRequest(req, LoginAccount)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": acc.Email,
		"iat":   time.Now().Unix(),
	})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("TestLoginAccount failed. Expected status %d - found %d", status, http.StatusOK)
	}

	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if signedToken, _ := token.SignedString([]byte(os.Getenv("MYSECRET"))); signedToken != result["token"] {
		t.Errorf("TestLoginAccount failed. JWT tokens differ. Generated %s - Received %s", signedToken, result["token"])
	}

}
