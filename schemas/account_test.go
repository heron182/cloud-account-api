package schemas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	InitDb()
	fmt.Println("Running main")

	status := m.Run()

	os.Exit(status)

	db.Database("accounts").Drop(context.Background())
}

// func TestCreateAccount(t *testing.T) {
// 	payload := []byte(`{"name": "John", "email: "john@example.com", "password": "123pass"}`)

// 	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(payload))
// 	rr := httptest.NewRecorder()

// 	handler := http.HandlerFunc(handlers.CreateAccount)
// 	handler.ServeHTTP(rr, req)
// }

func TestCreateAccount(t *testing.T) {
	acc := Account{Name: "John", Email: "john@example.com", Password: "123pass"}
	_, err := acc.Create()

	if err != nil {
		t.Errorf("Failed %s", err)
	}

	_, err = db.Database("accounts").
		Collection("accounts").
		FindOne(context.Background(), bson.M{"_id": acc.ID}).DecodeBytes()

	if err != nil {
		t.Errorf("Failed %s", err)
	}
}
