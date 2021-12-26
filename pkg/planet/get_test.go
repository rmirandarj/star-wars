package planet

import (
	"context"
	"errors"
	"testing"
	"time"

	"star-wars/pkg/testutils/docker"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_service_GetByID(t *testing.T) {
	mongoServer := docker.NewMongo()
	mongoServer.WithTestPort(t).
		Start(t)
	defer mongoServer.Stop()

	mongo := mongoCollection(mongoServer.GetHost())

	s := NewService(mongo, 2*time.Second)
	id, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")

	existingPlanet := Planet{
		ID:   id,
		Name: "Mars",
	}

	ctx := context.Background()
	_, err := s.db.InsertOne(ctx, existingPlanet)
	if err != nil {
		t.Fatalf("service.Insert() an error occurred inserting a planet for test")
	}

	tests := []struct {
		name        string
		givenID     string
		wantErr     bool
		wantErrType error
	}{
		{
			name:    "when the planet exists, then it should retrieve the planet with success",
			givenID: existingPlanet.ID.Hex(),
			wantErr: false,
		},
		{
			name:        "when the planet does not exist, then it should nothing",
			givenID:     "5f165e2e4de9b442e60b3905",
			wantErr:     true,
			wantErrType: ErrPlanetNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetByID(ctx, tt.givenID)
			if err != nil || tt.wantErr {
				if errors.Is(err, tt.wantErrType) {
					return
				}
				t.Fatalf("service.GetByID() errType = %v, wantErrType %v", err, tt.wantErrType)
			}
			assert.Equal(t, tt.givenID, got.ID.Hex(), "service.GetByID() unexpected planet id")
		})
	}
}
