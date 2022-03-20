package temperatures

import (
	"backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var temperaturesDb *utils.MongoInstance
var mongoCollection *mongo.Collection

func TemperaturesApp(app *fiber.App) {
	AddRoutes(app)
	ConnectToMongo()
}

func AddRoutes(app *fiber.App) {
	app.Get("/temperatures", ListTemperatures)
	app.Post("/temperatures", CreateTemperature)
}

func ConnectToMongo() {
	mongoInstance, err := utils.ConnectMongo()
	if err != nil {
		panic(err)
	}
	temperaturesDb = mongoInstance
	mongoCollection = temperaturesDb.Db.Collection("temperatures")
}
