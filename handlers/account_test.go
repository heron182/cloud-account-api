package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

func TestCreateAccount(t *testing.T) {
	defer tearDown()
	payload := []byte(`{"name": "John", "email": "john@createme.com", "password": "123pass"}`)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateAccount)
	handler.ServeHTTP(rr, req)

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
	acc.Create() // Password is mutated to a hash

	acc.Password = plainPassword
	if err := acc.CheckCredentials(); err != nil {
		t.Errorf("TestLoginAccount failed - %s", err)
	}

}
