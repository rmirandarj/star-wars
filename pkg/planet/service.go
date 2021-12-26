package planet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Planet struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Service struct {
	db *mongo.Collection
	timeout time.Duration
}

func NewService(db *mongo.Collection, timeout time.Duration) *Service {
	return &Service{
		db: db,
		timeout: timeout,
	}
}
