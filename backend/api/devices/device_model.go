package devices

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Device struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	SensorID    string `json:"sensorId" bson:"sensor_id" validate:"required"`
	Name        string `json:"name" bson:"name" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
	SensorType  string `json:"sensorType" bson:"sensorType" validate:"required"`
	Unit        string `json:"unit" bson:"unit" validate:"required"`
}

func (t *Device) GetId() string {
	return t.ID
}

func (t *Device) SetId(id primitive.ObjectID) {
	t.ID = id.Hex()
}

func (t Device) GetCollection() *mongo.Collection {
	return devicesCollection
}

func (t Device) Validate() error {
	return validator.New().Struct(t)
}
