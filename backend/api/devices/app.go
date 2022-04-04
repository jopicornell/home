package devices

import (
	"backend/lib"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var homeMongo *lib.MongoInstance
var devicesCollection *mongo.Collection

func App() *fiber.App {
	app := fiber.New()
	app.Get("/devices", ListDevicesController)
	app.Post("/devices", CreateDeviceController)
	ConnectToMongo()
	return app
}

func ConnectToMongo() {
	mongoInstance, err := lib.ConnectMongo()
	if err != nil {
		panic(err)
	}
	homeMongo = mongoInstance
	devicesCollection = homeMongo.Db.Collection("devices")
}
