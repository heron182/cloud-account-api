package schemas

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Account struct
type Account struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

// Create Creates an Account
func (acc *Account) Create() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	acc.ID = primitive.NewObjectID()

	if err := acc.hashPassword(); err != nil {
		log.Fatal(err)
	}

	_, err := accountCollection.InsertOne(ctx, acc)
	if err != nil {
		return err
	}

	return nil
}

// Get an account by ID
func (acc *Account) Get(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return accountCollection.FindOne(ctx, bson.M{"_id": id}).Decode(acc)
}

// CheckCredentials check if provided credentials exists in db
func (acc *Account) CheckCredentials() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	plainPassword := acc.Password

	if err := accountCollection.FindOne(ctx, bson.M{"email": acc.Email}).Decode(acc); err != nil {
		return errors.New("Invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(plainPassword)); err != nil {
		return errors.New("Invalid password")
	}

	return nil
}

// HashPassword hashes a plain password and assigns to Password
func (acc *Account) hashPassword() error {
	cost := 5
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(acc.Password), cost)
	acc.Password = string(hashedPwd)

	return err
}
