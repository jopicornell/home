package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func ConnectMongo() (*MongoInstance, error) {
	initDBUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	initDBPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	initDBName := os.Getenv("MONGO_INITDB_DATABASE")
	dbHost := os.Getenv("MONGO_DB_HOST")
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:27017/", initDBUsername, initDBPassword, dbHost)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(initDBName)

	if err != nil {
		return nil, err
	}

	mg := &MongoInstance{
		Client: client,
		Db:     db,
	}

	return mg, nil
}
