package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/util/enum"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pokemons = []model.Pokemon{
	{ID: 1, Name: "bulbasaur"},
	{ID: 2, Name: "ivysaur"},
	{ID: 3, Name: "venusaur"},
}

type mockUseCase struct {
	mock.Mock
}

func (muc mockUseCase) GetAllPokemons() []model.Pokemon {
	arg := muc.Called()
	return arg.Get(0).([]model.Pokemon)
}

func (muc mockUseCase) GetPokemonById(id int) (*model.Pokemon, error) {
	arg := muc.Called()
	return arg.Get(0).(*model.Pokemon), arg.Error(1)
}

func (muc mockUseCase) SetPokemons(pokemons []model.Pokemon) {
	// Just don't fail :)
}

func (muc mockUseCase) FetchPokemonsFromApi() ([]model.Pokemon, error) {
	arg := muc.Called()
	return arg.Get(0).([]model.Pokemon), arg.Error(1)
}

func (muc mockUseCase) FilterPokemonsConcurrently(enum.OddEven, int, int, int) []model.Pokemon {
	arg := muc.Called()
	return arg.Get(0).([]model.Pokemon)
}

// TestController_GetAllPokemons - Test controller get all pokemons
func TestController_GetAllPokemons(t *testing.T) {
	testCases := []struct {
		name             string
		useCasePokemons  []model.Pokemon
		expectedResponse []model.Pokemon
	}{
		{
			name:             "get all pokemons",
			useCasePokemons:  pokemons,
			expectedResponse: pokemons,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		muc := mockUseCase{}
		muc.On("GetAllPokemons").Return(tc.useCasePokemons)

		ctl := New(muc)

		r.GET("/pokemons", ctl.GetAllPokemons)

		c.Request, _ = http.NewRequest(http.MethodGet, "/pokemons", bytes.NewBuffer([]byte("{}")))

		r.ServeHTTP(w, c.Request)

		b, _ := ioutil.ReadAll(w.Body)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		} else {
			var response []model.Pokemon
			json.Unmarshal(b, &response)
			assert.Equal(t, tc.expectedResponse, response)
		}
	}
}

type controllerResponse struct {
	Message string `json:"message"`
}

// TestController_GetPokemonById - Test controller GetPokemonById
func TestController_GetPokemonById(t *testing.T) {
	testCases := []struct {
		name              string
		id                string
		useCasePokemon    *model.Pokemon
		expectedResponse  *model.Pokemon
		expectedErrorCode int
		error             error
	}{
		{
			name:              "get all pokemons",
			id:                "1",
			useCasePokemon:    &pokemons[0],
			expectedResponse:  &pokemons[0],
			expectedErrorCode: 0,
			error:             nil,
		},
		{
			name:              "get unparseable id pokemon",
			id:                "non-int",
			useCasePokemon:    &pokemons[0],
			expectedResponse:  nil,
			expectedErrorCode: http.StatusBadRequest,
			error:             errors.New("Failed to convert id non-int to int"),
		},
		{
			name:              "get wrong id pokemon",
			id:                "4",
			useCasePokemon:    nil,
			expectedResponse:  nil,
			expectedErrorCode: http.StatusNotFound,
			error:             errors.New("UseCase thrown error"),
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		muc := mockUseCase{}
		muc.On("GetPokemonById").Return(tc.useCasePokemon, tc.error)

		ctl := New(muc)

		r.GET("/pokemons/:id", ctl.GetPokemonById)

		c.Request, _ = http.NewRequest(http.MethodGet, "/pokemons/"+tc.id, bytes.NewBuffer([]byte("{}")))

		r.ServeHTTP(w, c.Request)

		b, _ := ioutil.ReadAll(w.Body)
		if w.Code != http.StatusOK {
			assert.Equal(t, tc.expectedErrorCode, w.Code)
			var cr controllerResponse
			json.Unmarshal(b, &cr)
			assert.Equal(t, tc.error.Error(), cr.Message)
		} else {
			var response model.Pokemon
			json.Unmarshal(b, &response)
			assert.Equal(t, tc.expectedResponse, &response)
		}
	}
}

// TestController_FetchPokemonsFromApi - Test controller FetchPokemonsFromApi method
func TestController_FetchPokemonsFromApi(t *testing.T) {
	testCases := []struct {
		name              string
		useCasePokemons   []model.Pokemon
		expectedResponse  []model.Pokemon
		expectedErrorCode int
		error             error
	}{
		{
			name:              "get all pokemons",
			useCasePokemons:   pokemons,
			expectedResponse:  pokemons,
			expectedErrorCode: 0,
			error:             nil,
		},
		{
			name:              "failed api request",
			useCasePokemons:   pokemons,
			expectedResponse:  []model.Pokemon{},
			expectedErrorCode: http.StatusInternalServerError,
			error:             errors.New("Internal Server Error"),
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		muc := mockUseCase{}
		muc.On("FetchPokemonsFromApi").Return(tc.useCasePokemons, tc.error)
		muc.On("GetAllPokemons").Return(tc.useCasePokemons)

		ctl := New(muc)

		r.GET("/pokemons/fetch", ctl.FetchPokemonsFromApi)

		c.Request, _ = http.NewRequest(http.MethodGet, "/pokemons/fetch", bytes.NewBuffer([]byte("{}")))

		r.ServeHTTP(w, c.Request)

		b, _ := ioutil.ReadAll(w.Body)
		if w.Code != http.StatusOK {
			assert.Equal(t, tc.expectedErrorCode, w.Code)
			var cr controllerResponse
			json.Unmarshal(b, &cr)
			assert.Equal(t, tc.error.Error(), cr.Message)
		} else {
			var response []model.Pokemon
			json.Unmarshal(b, &response)
			assert.Equal(t, tc.expectedResponse, response)
		}
	}
}

// TestController_FilterPokemonsConcurrently - Test controller FilterPokemonsConcurrently method
func TestController_FilterPokemonsConcurrently(t *testing.T) {
	testCases := []struct {
		name              string
		query             string
		useCasePokemons   []model.Pokemon
		expectedResponse  []model.Pokemon
		expectedErrorCode int
		error             error
	}{
		{
			name:              "default values",
			query:             "type=odd",
			useCasePokemons:   pokemons,
			expectedResponse:  pokemons,
			expectedErrorCode: 0,
			error:             nil,
		},
		{
			name:              "missing required value",
			query:             "",
			useCasePokemons:   pokemons,
			expectedResponse:  nil,
			expectedErrorCode: http.StatusBadRequest,
			error:             errors.New("'type' param error.  is not a valid input. Must be either 'odd' or 'even'"),
		},
		{
			name:              "wrong required value",
			query:             "type=par",
			useCasePokemons:   pokemons,
			expectedResponse:  nil,
			expectedErrorCode: http.StatusBadRequest,
			error:             errors.New("'type' param error. par is not a valid input. Must be either 'odd' or 'even'"),
		},
		{
			name:              "wrong items value",
			query:             "type=even&items=items",
			useCasePokemons:   pokemons,
			expectedResponse:  nil,
			expectedErrorCode: http.StatusBadRequest,
			error:             errors.New("'items' param error. Cannot convert items to int"),
		},
		{
			name:              "wrong items_per_workers value",
			query:             "type=even&items=4&items_per_workers=ipw",
			useCasePokemons:   pokemons,
			expectedResponse:  nil,
			expectedErrorCode: http.StatusBadRequest,
			error:             errors.New("'items_per_workers' param error. Cannot convert ipw to int"),
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		muc := mockUseCase{}
		muc.On("FilterPokemonsConcurrently").Return(tc.useCasePokemons)

		ctl := New(muc)

		r.GET("/pokemons/filter", ctl.FilterPokemonsConcurrently)

		c.Request, _ = http.NewRequest(http.MethodGet, "/pokemons/filter?"+tc.query, bytes.NewBuffer([]byte("{}")))

		r.ServeHTTP(w, c.Request)

		b, _ := ioutil.ReadAll(w.Body)
		if w.Code != http.StatusOK {
			assert.Equal(t, tc.expectedErrorCode, w.Code)
			var cr controllerResponse
			json.Unmarshal(b, &cr)
			assert.Equal(t, tc.error.Error(), cr.Message)
		} else {
			var response []model.Pokemon
			json.Unmarshal(b, &response)
			assert.Equal(t, tc.expectedResponse, response)
		}
	}
}
