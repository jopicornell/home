package lib

import (
	"backend/log"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func InsertOne(model MongoModel) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := model.GetCollection().InsertOne(ctx, model)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func InsertFromBodyAndGetId(ctx *fiber.Ctx, model MongoModel) (primitive.ObjectID, error) {
	logger := log.Logger.WithField("requestId", ctx.Locals("request_id"))
	err := ctx.BodyParser(model)
	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return primitive.NilObjectID, err
	}
	err = model.Validate()
	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return primitive.NilObjectID, err
	}
	result, err := model.GetCollection().InsertOne(ctx.Context(), model)
	if err != nil {
		logger.Error(err)
		_ = ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting document",
		})
		return primitive.NilObjectID, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		_ = ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting inserted id",
		})
		logger.Error(errors.New("Error getting inserted id"))
		return primitive.NilObjectID, err
	}
	model.SetId(oid)
	logrus.Info("Inserted document with id: ", oid.Hex())
	ctx.Append("Location", fmt.Sprintf("%s/%s", ctx.Path(), oid.Hex()))
	return oid, nil
}

func InsertFromBody(ctx *fiber.Ctx, model MongoModel) error {
	_, err := InsertFromBodyAndGetId(ctx, model)
	return err
}
