package planet

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrPlanetNotFound = errors.New("planet not found")
)

func (s *Service) GetByID(ctx context.Context, id string) (Planet, error) {
	var planet Planet

	objectID, _ := primitive.ObjectIDFromHex(id)
	result := s.db.FindOne(ctx, bson.M{"_id": objectID})
	err := result.Decode(&planet)

	if errors.Is(err, driver.ErrNoDocuments) {
		return planet, ErrPlanetNotFound
	}
		
	return planet, err
}
