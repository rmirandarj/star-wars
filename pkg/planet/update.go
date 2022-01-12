package planet

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) Update(ctx context.Context, planetDocument Planet) (int64, error) {

	filter := bson.M{"_id": planetDocument.ID}
	update := bson.M{"$set": bson.M{
		"name": planetDocument.Name,
	}}

	result, err := s.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return 0, err
	} else if result.MatchedCount == 0 {
		return 0, ErrPlanetNotFound
	}

	return result.MatchedCount, err
}
