package device_values

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeviceValue struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Value  float64 `json:"value" bson:"value"`
	Sensor string  `json:"sensor" bson:"sensor"`
}

func (t *DeviceValue) GetId() string {
	return t.ID
}

func (t *DeviceValue) SetId(id primitive.ObjectID) {
	t.ID = id.Hex()
}

func (t DeviceValue) GetCollection() *mongo.Collection {
	return deviceValuesCollection
}

func (t DeviceValue) Validate() error {
	return validator.New().Struct(t)
}
