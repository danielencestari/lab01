package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-cep-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do CEP Service
type MockCEPService struct {
	mock.Mock
}

func (m *MockCEPService) GetLocationByCEP(cep string) (*models.LocationInfo, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationInfo), args.Error(1)
}

// Mock do Weather Service
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetTemperatureByCity(city, state string) (*models.TemperatureResponse, error) {
	args := m.Called(city, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TemperatureResponse), args.Error(1)
}

func setupRouter(handler *WeatherHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", handler.HealthCheck)
	router.GET("/temperature/:cep", handler.GetTemperatureByCEP)
	return router
}

func TestWeatherHandler_HealthCheck(t *testing.T) {
	// Setup
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Weather CEP API is running", response["message"])
}

func TestWeatherHandler_GetTemperatureByCEP_Success(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Mock responses
	locationInfo := &models.LocationInfo{
		City:  "São Paulo",
		State: "SP",
		CEP:   "01310-100",
	}
	
	tempResponse := &models.TemperatureResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.0,
	}

	// Setup mock expectations
	mockCEPService.On("GetLocationByCEP", "01310100").Return(locationInfo, nil)
	mockWeatherService.On("GetTemperatureByCity", "São Paulo", "SP").Return(tempResponse, nil)

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.TemperatureResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 25.0, response.TempC)
	assert.Equal(t, 77.0, response.TempF)
	assert.Equal(t, 298.0, response.TempK)

	// Verify mock calls
	mockCEPService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}

func TestWeatherHandler_GetTemperatureByCEP_InvalidCEP(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Setup mock expectations
	mockCEPService.On("GetLocationByCEP", "123").Return(nil, errors.New("invalid zipcode"))

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/temperature/123", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "invalid zipcode", w.Body.String())

	// Verify mock calls
	mockCEPService.AssertExpectations(t)
}

func TestWeatherHandler_GetTemperatureByCEP_CEPNotFound(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Setup mock expectations
	mockCEPService.On("GetLocationByCEP", "99999999").Return(nil, errors.New("can not find zipcode"))

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/temperature/99999999", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "can not find zipcode", w.Body.String())

	// Verify mock calls
	mockCEPService.AssertExpectations(t)
}

func TestWeatherHandler_GetTemperatureByCEP_WeatherServiceError(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Mock responses
	locationInfo := &models.LocationInfo{
		City:  "São Paulo",
		State: "SP",
		CEP:   "01310-100",
	}

	// Setup mock expectations
	mockCEPService.On("GetLocationByCEP", "01310100").Return(locationInfo, nil)
	mockWeatherService.On("GetTemperatureByCity", "São Paulo", "SP").Return(nil, errors.New("weather API error"))

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error fetching weather data", w.Body.String())

	// Verify mock calls
	mockCEPService.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
}

func TestWeatherHandler_GetTemperatureByCEP_EmptyCEP(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request with empty CEP
	req, _ := http.NewRequest("GET", "/temperature/", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions - Gin should return 404 for missing parameter
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestWeatherHandler_GetTemperatureByCEP_CEPServiceInternalError(t *testing.T) {
	// Setup mocks
	mockCEPService := new(MockCEPService)
	mockWeatherService := new(MockWeatherService)

	// Setup mock expectations - internal error different from validation errors
	mockCEPService.On("GetLocationByCEP", "01310100").Return(nil, errors.New("internal database error"))

	// Setup handler and router
	handler := NewWeatherHandler(mockCEPService, mockWeatherService)
	router := setupRouter(handler)

	// Create request
	req, _ := http.NewRequest("GET", "/temperature/01310100", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "internal server error", w.Body.String())

	// Verify mock calls
	mockCEPService.AssertExpectations(t)
} 