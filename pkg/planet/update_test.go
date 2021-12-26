package planet

import (
	"context"
	"testing"
	"time"

	"star-wars/pkg/testutils/docker"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_service_Update(t *testing.T) {
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
		t.Fatalf("service.Update() an error occurred inserting a planet for test")
	}

	id2, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3905")

	tests := []struct {
		name             string
		givenPlanet      Planet
		wantMatchedCount int64
		wantError        bool
	}{
		{
			name: "when the planet exists, then it should update with success",
			givenPlanet: Planet{
				ID:   existingPlanet.ID,
				Name: "New Mars",
			},
			wantMatchedCount: 1,
		},
		{
			name: "when the planet does not exist, then it should return an error",
			givenPlanet: Planet{
				ID:   id2,
				Name: "New Mars",
			},
			wantMatchedCount: 0,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchedCount, err := s.Update(ctx, tt.givenPlanet)
			if err != nil {
				t.Fatalf("service.Update() unexpected error %v", err)
			}
			assert.Equal(t, tt.wantMatchedCount, matchedCount, "service.Update() unexpected planet matchedCount")

			if matchedCount > 0 {
				p, err := s.GetByID(ctx, tt.givenPlanet.ID.Hex())
				if err != nil {
					t.Fatalf("service.Update() an error occurred retrieving a planet for test")
				}
				assert.Equal(t, tt.givenPlanet.Name, p.Name, "service.Update() unexpected planet name")
			}
		})
	}
}

func mongoCollection(host string) *driver.Collection {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	client, err := driver.Connect(ctx, options.Client().ApplyURI(host))

	if err != nil {
		log.Fatal("Error trying to connect to the database")
	}

	return client.Database("planet").Collection("planet")
}
