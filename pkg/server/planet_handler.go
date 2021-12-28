package server

import (
	"context"
	"errors"
	"net/http"

	"star-wars/pkg/planet"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanetDTO struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

const (
	AppName        = "communications-logs-api"
	XApplicationId = "x-application-id"
	XWorkspaceId   = "x-workspace-id"
)

type PlanetInserter interface {
	Insert(ctx context.Context, planet planet.Planet) (planet.Planet, error)
}

func (a *App) handleCreatePlanet(saver PlanetInserter) http.Handler {
	type planetRequest struct {
		Name string `json:"name" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := loggerFromRequest(r)

		var planetRequest planetRequest

		if err := decodeAndValidate(w, r, &planetRequest); err != nil {
			logger.Error(err.Error())
			return
		}

		doc := planet.Planet{
			Name: planetRequest.Name,
		}
		saved, err := saver.Insert(ctx, doc)
		if err != nil {
			logger.Error(err.Error())
			writeJsonResponse(w, http.StatusInternalServerError, errorMessage{Message: "failed to insert the planet", ErrorCode: "WA:002"})
			return
		}

		dto := PlanetDTO{
			ID:   saved.ID,
			Name: saved.Name,
		}
		writeJsonResponse(w, http.StatusCreated, dto)
	})
}

type PlanetUpdater interface {
	Update(ctx context.Context, planetDocument planet.Planet) (int64, error)
}

func (a *App) handleUpdatePlanet(planetUpdater PlanetUpdater) http.HandlerFunc {
	type planetRequest struct {
		Name string `json:"name" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := loggerFromRequest(r)
		id := mux.Vars(r)["id"]

		objectID, _ := primitive.ObjectIDFromHex(id)

		var planetRequest planetRequest

		if err := decodeAndValidate(w, r, &planetRequest); err != nil {
			logger.Error(err.Error())
			return
		}

		doc := planet.Planet{
			ID:   objectID,
			Name: planetRequest.Name,
		}
		matchedCount, err := planetUpdater.Update(ctx, doc)
		if err != nil {
			logger.Error(err.Error())
			writeJsonResponse(w, http.StatusInternalServerError, errorMessage{Message: "failed to update the planet", ErrorCode: "WA:001"})
			return
		}

		if matchedCount == 0 {
			logger.Error(err.Error())
			writeJsonResponse(w, http.StatusNotFound, errorMessage{Message: "planet not found", ErrorCode: "WA:003"})
			return
		}
		writeJsonResponse(w, http.StatusNoContent, nil)
	})
}

type PlanetGetter interface {
	GetByID(context.Context, string) (planet.Planet, error)
}

func (a *App) handleGetPlanetByID(planetGetter PlanetGetter) http.HandlerFunc {
	type planetDTO struct {
		ID   primitive.ObjectID `json:"id"`
		Name string             `json:"name"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := mux.Vars(r)["id"]

		got, err := planetGetter.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, planet.ErrPlanetNotFound) {
				writeJsonResponse(rw, http.StatusNotFound, errorMessage{ErrorCode: "WA:003", Message: "planet not found"})
				return
			}
			writeJsonResponse(rw, http.StatusInternalServerError, errorMessage{ErrorCode: "WA:004", Message: "failed to retrieve a planet by id"})
			return
		}

		writeJsonResponse(rw, http.StatusOK, planetDTO{
			ID:   got.ID,
			Name: got.Name,
		})
	}
}
