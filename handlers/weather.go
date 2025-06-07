package handlers

import (
	"net/http"
	"strings"
	"weather-cep-api/services"

	"github.com/gin-gonic/gin"
)

// WeatherHandler gerencia requisições relacionadas ao clima
type WeatherHandler struct {
	cepService     services.CEPServiceInterface
	weatherService services.WeatherServiceInterface
}

// NewWeatherHandler cria uma nova instância do handler de clima
func NewWeatherHandler(cepService services.CEPServiceInterface, weatherService services.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		cepService:     cepService,
		weatherService: weatherService,
	}
}

// GetTemperatureByCEP busca temperatura por CEP
// GET /temperature/:cep
func (h *WeatherHandler) GetTemperatureByCEP(c *gin.Context) {
	// Extrai CEP dos parâmetros da URL
	cep := c.Param("cep")
	if cep == "" {
		c.String(http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	// 1. Busca informações de localização pelo CEP
	location, err := h.cepService.GetLocationByCEP(cep)
	if err != nil {
		if strings.Contains(err.Error(), "invalid zipcode") {
			c.String(http.StatusUnprocessableEntity, "invalid zipcode")
			return
		}
		if strings.Contains(err.Error(), "can not find zipcode") {
			c.String(http.StatusNotFound, "can not find zipcode")
			return
		}
		// Erro interno do servidor
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	// 2. Busca temperatura pela cidade/estado
	temperature, err := h.weatherService.GetTemperatureByCity(location.City, location.State)
	if err != nil {
		// Log do erro para debug (em produção usar logger)
		c.String(http.StatusInternalServerError, "error fetching weather data")
		return
	}

	// 3. Retorna resposta com temperaturas
	c.JSON(http.StatusOK, temperature)
}

// HealthCheck endpoint para verificação de saúde da aplicação
// GET /health
func (h *WeatherHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Weather CEP API is running",
	})
} 