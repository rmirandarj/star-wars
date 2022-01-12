package server

import "star-wars/pkg/planet"

type container struct {
	planetInserter PlanetInserter
	planetUpdater  PlanetUpdater
	planetGetter   PlanetGetter
}

func NewContainer(planetService *planet.Service) *container {
	return &container{
		planetInserter: planetService,
		planetUpdater:  planetService,
		planetGetter:   planetService,
	}
}
