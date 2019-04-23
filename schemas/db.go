package schemas

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	databaseName          string = "accounts"
	accountCollectionName string = "accounts"
)

// Db is the main database interface
var Db *mongo.Client
var accountCollection *mongo.Collection

// InitDb starts db connections
func InitDb(databaseURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Db, err = mongo.Connect(ctx, options.Client().ApplyURI(databaseURI))
	accountCollection = Db.Database(databaseName).Collection(accountCollectionName)

	if err != nil {
		log.Fatal(err)
	}

	if err = Db.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

}
