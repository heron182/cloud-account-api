package schemas

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("-- No .env file encountered in the root project folder")
	}

	InitDb(os.Getenv("DATABASE_URI"))
}

func tearDown() {
	Db.Database("accounts").Drop(context.Background())
}

func TestMain(m *testing.M) {
	setup()

	status := m.Run()

	tearDown()

	os.Exit(status)

}

func TestCreate(t *testing.T) {
	defer tearDown()

	acc := Account{Name: "John", Email: "john@example.com", Password: "123pass"}
	err := acc.Create()

	if err != nil {
		t.Errorf("TestCreate failed - %s", err)
	}

	err = accountCollection.
		FindOne(context.Background(), bson.M{"_id": acc.ID}).
		Decode(&acc)

	if err != nil {
		t.Errorf("TestCreateAccount failed - %s", err)
	}
}

func TestHashPassword(t *testing.T) {
	defer tearDown()

	plainPassword := "123password"
	acc := Account{Password: plainPassword}
	acc.hashPassword()

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(plainPassword)); err != nil {
		t.Errorf("TestHashPassword failed - %s", err)
	}

}

func TestCheckCredentials(t *testing.T) {
	defer tearDown()

	plainPassword := "123password"
	acc := Account{Name: "John", Email: "john@example.com", Password: plainPassword}
	acc.Create()

	acc.Password = plainPassword
	if err := acc.CheckCredentials(); err != nil {
		t.Errorf("TestCheckCredentials failed - %s - %s", err, acc.Password)
	}

}
