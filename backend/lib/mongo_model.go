package lib

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoModel interface {
	GetCollection() *mongo.Collection
	GetId() string
	SetId(primitive.ObjectID)
	Validate() error
}
