package server

import "star-wars/pkg/planet"

type Container struct {
	planetInserter PlanetInserter
	planetUpdater  PlanetUpdater
	planetGetter   PlanetGetter
}

func NewContainer(planetService *planet.Service) *Container {
	return &Container{
		planetInserter: planetService,
		planetUpdater:  planetService,
		planetGetter:   planetService,
	}
}
