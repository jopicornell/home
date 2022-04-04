package device_values

import (
	"backend/lib"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var homeMongo *lib.MongoInstance
var deviceValuesCollection *mongo.Collection

func App() *fiber.App {
	app := fiber.New()
	AddRoutes(app)
	ConnectToMongo()
	return app
}

func AddRoutes(app *fiber.App) {
	app.Get("/devices/:device_id/values", ListDeviceValuesController)
	app.Post("/devices/:device_id/values", CreateDeviceValueController)
}

func ConnectToMongo() {
	mongoInstance, err := lib.ConnectMongo()
	if err != nil {
		panic(err)
	}
	homeMongo = mongoInstance
	deviceValuesCollection = homeMongo.Db.Collection("device_values")
}
