package docker

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"star-wars/pkg/testutils"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dockerMongo struct {
	*dockerBase
}

func NewMongo() Docker {
	base := &dockerBase{}
	dk := &dockerMongo{
		dockerBase: base,
	}
	base.adapter = dk
	return base
}

func (dm *dockerMongo) WithTestPort(t *testing.T) Docker {
	return dm.WithPort(testutils.FindAvailablePort("localhost"))
}

func (dm *dockerMongo) Wait() error {
	host := dm.GetHost()

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	client, err := driver.Connect(ctx, options.Client().ApplyURI(host))

	if err != nil {
		log.Fatal("Error trying to connect to the database.", err)
	}
	defer client.Disconnect(ctx)

	return err
}

func (dm *dockerMongo) GetHost() string {
	return fmt.Sprintf("mongodb://localhost:%s/planet?readPreference=primary", dm.port)
}

func (dm *dockerMongo) CreateImage() (pool *dockertest.Pool, resource *dockertest.Resource, err error) {
	pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	opts := dockertest.RunOptions{
		Repository: "mongo",
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("27017/tcp"): {
				{HostIP: "0.0.0.0", HostPort: dm.port},
			},
		},
		ExposedPorts: []string{dm.port},
		Env:          []string{},
	}

	resource, err = pool.RunWithOptions(&opts)

	if err != nil {
		return nil, nil, err
	}
	return
}
