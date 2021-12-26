package planet

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) Insert(ctx context.Context, planetDocument Planet) (Planet, error) {

	planetDocument.ID = primitive.NewObjectID()
	_, err := s.db.InsertOne(ctx, planetDocument)

	if err != nil {
		return planetDocument, err
	}

	return planetDocument, err
}
