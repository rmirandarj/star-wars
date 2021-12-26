package planet

import (
	"context"
	"testing"
	"time"

	"star-wars/pkg/testutils/docker"

	"github.com/stretchr/testify/assert"
)

func Test_service_Insert(t *testing.T) {
	mongoServer := docker.NewMongo()
	mongoServer.WithTestPort(t).
		Start(t)
	defer mongoServer.Stop()

	mongo := mongoCollection(mongoServer.GetHost())
	s := NewService(mongo, 2*time.Second)
	ctx := context.Background()

	tests := []struct {
		name        string
		givenPlanet Planet
		wantPlanet  Planet
	}{
		{
			name: "when given new planet, then it should insert with success",
			givenPlanet: Planet{
				Name: "Pluto",
			},
			wantPlanet: Planet{
				Name: "Pluto",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planet, err := s.Insert(ctx, tt.givenPlanet)
			if err != nil {
				t.Fatalf("service.Insert() unexpected error %v", err)
			}
			assert.Equal(t, tt.givenPlanet.Name, planet.Name, "service.Insert() unexpected planet name")
		})
	}
}
