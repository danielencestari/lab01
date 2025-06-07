package main

import (
	"log"
	"os"
	"weather-cep-api/handlers"
	"weather-cep-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis do arquivo .env (ignora erro se arquivo não existir)
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Arquivo .env não encontrado ou erro ao carregar: %v", err)
		log.Printf("Usando variáveis de ambiente do sistema")
	} else {
		log.Printf("Arquivo .env carregado com sucesso")
	}

	// Configura modo do Gin baseado na variável de ambiente
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Cria instâncias dos serviços
	cepService := services.NewCEPService()
	weatherService := services.NewWeatherService()

	// Cria instância do handler
	weatherHandler := handlers.NewWeatherHandler(cepService, weatherService)

	// Configura o router Gin
	router := gin.Default()

	// Adiciona middleware de CORS para permitir requisições de diferentes origens
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define as rotas
	router.GET("/health", weatherHandler.HealthCheck)
	router.GET("/temperature/:cep", weatherHandler.GetTemperatureByCEP)

	// Define a porta do servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicia o servidor
	log.Printf("Servidor iniciando na porta %s", port)
	log.Printf("Endpoints disponíveis:")
	log.Printf("  GET /health - Health check")
	log.Printf("  GET /temperature/:cep - Consulta temperatura por CEP")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
} 