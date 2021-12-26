package docker

import (
	"testing"

	"github.com/ory/dockertest/v3"
)

type Docker interface {
	CreateImage() (pool *dockertest.Pool, resource *dockertest.Resource, err error)
	Start(t *testing.T)
	Wait() error
	Stop() error
	GetHost() string
	GetPort() string
	WithTestPort(t *testing.T) Docker
	WithPort(port string) Docker
}

type dockerBase struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	port     string
	adapter  Docker
}

func (d *dockerBase) WithPort(port string) Docker {
	d.port = port
	return d
}

func (d *dockerBase) WithTestPort(t *testing.T) Docker {
	return d.adapter.WithTestPort(t)
}

func (d *dockerBase) Start(t *testing.T) {
	if d.port == "" {
		t.Fatalf("please inform a port to the docker image, ex: DockerKafka{}.WithTestPort(t)")
	}

	pool, resource, err := d.CreateImage()
	if err != nil {
		t.Fatalf("failed to create docker test container %s", err)
	}

	d.pool = pool
	d.resource = resource

	if err := d.Wait(); err != nil {
		d.Stop()
		t.Fatalf("container didn't start on time %s", err)
	}

}

func (d *dockerBase) Wait() error {
	return d.pool.Retry(d.adapter.Wait)
}

func (d *dockerBase) Stop() error {
	return d.pool.Purge(d.resource)
}

func (d *dockerBase) GetHost() string {
	return d.adapter.GetHost()
}

func (d *dockerBase) GetPort() string {
	return d.port
}

func (d *dockerBase) CreateImage() (pool *dockertest.Pool, resource *dockertest.Resource, err error) {
	return d.adapter.CreateImage()
}
