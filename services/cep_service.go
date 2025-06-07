package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-cep-api/models"
	"weather-cep-api/utils"
)

// CEPServiceInterface define o contrato para serviços de CEP
type CEPServiceInterface interface {
	GetLocationByCEP(cep string) (*models.LocationInfo, error)
}

// CEPService implementa o serviço de consulta de CEP
type CEPService struct {
	httpClient HTTPClientInterface
}

// NewCEPService cria uma nova instância do serviço de CEP
func NewCEPService() *CEPService {
	return &CEPService{
		httpClient: &http.Client{},
	}
}

// NewCEPServiceWithClient cria uma nova instância com HTTP client customizado (para testes)
func NewCEPServiceWithClient(client HTTPClientInterface) *CEPService {
	return &CEPService{
		httpClient: client,
	}
}

// GetLocationByCEP consulta informações de localização por CEP usando ViaCEP
func (s *CEPService) GetLocationByCEP(cep string) (*models.LocationInfo, error) {
	// Valida formato do CEP
	if !utils.IsValidCEP(cep) {
		return nil, fmt.Errorf("invalid zipcode")
	}

	// Normaliza o CEP (remove hífen)
	normalizedCEP := utils.NormalizeCEP(cep)

	// Constrói URL da ViaCEP
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", normalizedCEP)

	// Faz a requisição HTTP
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching CEP data: %w", err)
	}
	defer resp.Body.Close()

	// Verifica status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can not find zipcode")
	}

	// Decodifica a resposta JSON
	var viaCEPResp models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResp); err != nil {
		return nil, fmt.Errorf("error decoding CEP response: %w", err)
	}

	// Verifica se o CEP foi encontrado (ViaCEP retorna erro: true quando não encontra)
	if viaCEPResp.Erro {
		return nil, fmt.Errorf("can not find zipcode")
	}

	// Valida se os campos essenciais estão preenchidos
	if viaCEPResp.Localidade == "" {
		return nil, fmt.Errorf("can not find zipcode")
	}

	// Converte para o modelo interno
	locationInfo := &models.LocationInfo{
		City:  viaCEPResp.Localidade,
		State: viaCEPResp.UF,
		CEP:   utils.FormatCEP(normalizedCEP),
	}

	return locationInfo, nil
} 