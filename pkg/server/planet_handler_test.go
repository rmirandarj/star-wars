package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"star-wars/pkg/planet"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type planetInserterMock struct {
	result planet.Planet
	err    error
}

func (a planetInserterMock) Insert(ctx context.Context, planet planet.Planet) (planet.Planet, error) {
	return a.result, a.err
}

func Test_handleInsertPlanet(t *testing.T) {
	t.Parallel()

	objectID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	tests := []struct {
		name               string
		givenBody          string
		planetInserterMock planetInserterMock
		wantStatusCode     int
		wantResponseBody   string
	}{
		{
			name:      "when payload is valid and no error then it should return 201 status",
			givenBody: `{"name": "Mars"}`,
			planetInserterMock: planetInserterMock{
				result: planet.Planet{ID: objectID, Name: "Mars"},
				err:    nil,
			},
			wantStatusCode:   201,
			wantResponseBody: `{"id":"5f165e2e4de9b442e60b3904","name":"Mars"}`,
		},
		{
			name:      "when payload is invalid then it should return 400 status",
			givenBody: `{: "Mars"}`,
			planetInserterMock: planetInserterMock{
				result: planet.Planet{},
				err:    nil,
			},
			wantStatusCode:   400,
			wantResponseBody: `{"error_code":"WA:007","message":"failed to decode payload"}`,
		},
		{
			name:      "when required field not sent then it should return 422 status",
			givenBody: `{"test": "Mars"}`,
			planetInserterMock: planetInserterMock{
				result: planet.Planet{},
				err:    nil,
			},
			wantStatusCode:   422,
			wantResponseBody: `{"error_code":"WA:001","message":"payload is invalid","details":[{"name":"Name","reason":"Key: 'planetRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}]}`,
		},
		{
			name:      "when payload is valid and can't save it then it should return 500 status",
			givenBody: `{"name": "Mars"}`,
			planetInserterMock: planetInserterMock{
				result: planet.Planet{},
				err:    errors.New("Database Error"),
			},
			wantStatusCode:   500,
			wantResponseBody: `{"error_code":"WA:002","message":"failed to insert the planet"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var app App
			app.container = &container{planetInserter: tc.planetInserterMock}
			app.RegisterRoutes()

			req, _ := http.NewRequest("POST", "/v1/planets", strings.NewReader(tc.givenBody))
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)
			if rr.Code != tc.wantStatusCode {
				t.Errorf("handleInsertPlanet() status code = %v, want %v", rr.Code, tc.wantStatusCode)
			}
			if got, _ := ioutil.ReadAll(rr.Result().Body); string(got) != tc.wantResponseBody {
				t.Errorf("handleInsertPlanet() body = %v, want %v", string(got), tc.wantResponseBody)
			}
		})
	}
}

type planetUpdaterMock struct {
	matchedCount int64
	err          error
}

func (a planetUpdaterMock) Update(ctx context.Context, planet planet.Planet) (int64, error) {
	return a.matchedCount, a.err
}

func Test_handleUpdatePlanet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		givenPlanetID     string
		givenBody         string
		planetUpdaterMock planetUpdaterMock
		wantStatusCode    int
		wantResponseBody  string
	}{
		{
			name:          "when payload is valid and no error then it should return 204 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			givenBody:     `{"name": "Mars"}`,
			planetUpdaterMock: planetUpdaterMock{
				matchedCount: 1,
				err:          nil,
			},
			wantStatusCode: 204,
		},
		{
			name:          "when payload is invalid then it should return 400 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			givenBody:     `{: "Mars"}`,
			planetUpdaterMock: planetUpdaterMock{
				err: nil,
			},
			wantStatusCode:   400,
			wantResponseBody: `{"error_code":"WA:007","message":"failed to decode payload"}`,
		},
		{
			name:          "When required field not sent then it should return 422 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			givenBody:     `{"test": "Mars"}`,
			planetUpdaterMock: planetUpdaterMock{
				err: nil,
			},
			wantStatusCode:   422,
			wantResponseBody: `{"error_code":"WA:001","message":"payload is invalid","details":[{"name":"Name","reason":"Key: 'planetRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}]}`,
		},
		{
			name:          "when the planet not found then it should return 404 status",
			givenPlanetID: "5f165e2e4de9b442e60b3905",
			givenBody:     `{"name": "Mars"}`,
			planetUpdaterMock: planetUpdaterMock{
				matchedCount: 0,
				err:          planet.ErrPlanetNotFound,
			},
			wantStatusCode:   404,
			wantResponseBody: `{"error_code":"WA:003","message":"planet not found"}`,
		},
		{
			name:          "when payload is valid and can't save it then it should return 500 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			givenBody:     `{"name": "Mars"}`,
			planetUpdaterMock: planetUpdaterMock{
				matchedCount: 0,
				err:          errors.New("Database Error"),
			},
			wantStatusCode:   500,
			wantResponseBody: `{"error_code":"WA:001","message":"failed to update the planet"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var app App
			app.container = &container{planetUpdater: tc.planetUpdaterMock}
			app.RegisterRoutes()

			req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/planets/%s", tc.givenPlanetID), strings.NewReader(tc.givenBody))
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)
			if rr.Code != tc.wantStatusCode {
				t.Errorf("handleUpdatePlanet() status code = %v, want %v", rr.Code, tc.wantStatusCode)
			}
			if got, _ := ioutil.ReadAll(rr.Result().Body); tc.wantResponseBody != "" && string(got) != tc.wantResponseBody {
				t.Errorf("handleUpdatePlanet() body = %v, want %v", string(got), tc.wantResponseBody)
			}
		})
	}
}

type planetGetterMock struct {
	result planet.Planet
	err    error
}

func (a planetGetterMock) GetByID(ctx context.Context, accountID string) (planet.Planet, error) {
	return a.result, a.err
}

func Test_handleGetAccount(t *testing.T) {
	t.Parallel()

	objectID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	tests := []struct {
		name             string
		givenPlanetID    string
		planetGetterMock planetGetterMock
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name:          "when planet id is informed and no error then it should return 200 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			planetGetterMock: planetGetterMock{
				result: planet.Planet{ID: objectID, Name: "Mars"},
				err:    nil,
			},
			wantStatusCode:   200,
			wantResponseBody: `{"id":"5f165e2e4de9b442e60b3904","name":"Mars"}`,
		},
		{
			name:          "when planet id is informed but got generic error from database then it should return 500 status",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			planetGetterMock: planetGetterMock{
				result: planet.Planet{ID: objectID, Name: "Mars"},
				err:    errors.New("database error"),
			},
			wantStatusCode:   500,
			wantResponseBody: `{"error_code":"WA:004","message":"failed to retrieve a planet by id"}`,
		},
		{
			name:          "when planet id is informed but couldn't find in database then it should return 404",
			givenPlanetID: "5f165e2e4de9b442e60b3904",
			planetGetterMock: planetGetterMock{
				result: planet.Planet{ID: objectID, Name: "Mars"},
				err:    planet.ErrPlanetNotFound,
			},
			wantStatusCode:   404,
			wantResponseBody: `{"error_code":"WA:003","message":"planet not found"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var app App
			app.container = &container{planetGetter: tc.planetGetterMock}
			app.RegisterRoutes()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/planets/%s", tc.givenPlanetID), nil)
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)
			if rr.Code != tc.wantStatusCode {
				t.Errorf("handleGetPlanet() status code = %v, want %v", rr.Code, tc.wantStatusCode)
			}
			if got, _ := ioutil.ReadAll(rr.Result().Body); tc.wantResponseBody != "" && string(got) != tc.wantResponseBody {
				t.Errorf("handleGetPlanet() body = %v, want %v", string(got), tc.wantResponseBody)
			}
		})
	}
}
