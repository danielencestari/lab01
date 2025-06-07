package services

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock do HTTP Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCEPService_GetLocationByCEP_Success(t *testing.T) {
	// Mock da resposta da ViaCEP
	mockResponse := `{
		"cep": "01310-100",
		"logradouro": "Avenida Paulista",
		"complemento": "",
		"bairro": "Bela Vista",
		"localidade": "São Paulo",
		"uf": "SP",
		"ibge": "3550308",
		"gia": "1004",
		"ddd": "11",
		"siafi": "7107"
	}`

	// Cria resposta HTTP mock
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponse)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	// Configura mock do HTTP client
	mockClient := new(MockHTTPClient)
	mockClient.On("Get", "https://viacep.com.br/ws/01310100/json/").Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewCEPServiceWithClient(mockClient)

	// Executa o teste
	result, err := service.GetLocationByCEP("01310-100")

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "São Paulo", result.City)
	assert.Equal(t, "SP", result.State)
	assert.Equal(t, "01310-100", result.CEP)

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestCEPService_GetLocationByCEP_InvalidCEP(t *testing.T) {
	service := NewCEPService()

	tests := []struct {
		name string
		cep  string
	}{
		{
			name: "CEP too short",
			cep:  "0131010",
		},
		{
			name: "CEP too long",
			cep:  "013101000",
		},
		{
			name: "CEP with letters",
			cep:  "01310a00",
		},
		{
			name: "Empty CEP",
			cep:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetLocationByCEP(tt.cep)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid zipcode")
		})
	}
}

func TestCEPService_GetLocationByCEP_NotFound(t *testing.T) {
	// Mock da resposta de CEP não encontrado
	mockResponse := `{
		"erro": true
	}`

	// Cria resposta HTTP mock
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponse)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	// Configura mock do HTTP client
	mockClient := new(MockHTTPClient)
	mockClient.On("Get", "https://viacep.com.br/ws/99999999/json/").Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewCEPServiceWithClient(mockClient)

	// Executa o teste
	result, err := service.GetLocationByCEP("99999-999")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "can not find zipcode")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestCEPService_GetLocationByCEP_ServerError(t *testing.T) {
	// Configura mock do HTTP client para retornar erro 500
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}

	mockClient := new(MockHTTPClient)
	mockClient.On("Get", "https://viacep.com.br/ws/01310100/json/").Return(resp, nil)

	// Cria service com HTTP client mockado
	service := NewCEPServiceWithClient(mockClient)

	// Executa o teste
	result, err := service.GetLocationByCEP("01310-100")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "can not find zipcode")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
}

func TestCEPService_GetLocationByCEP_HTTPError(t *testing.T) {
	// Configura mock do HTTP client para retornar erro de conexão
	mockClient := new(MockHTTPClient)
	mockClient.On("Get", "https://viacep.com.br/ws/01310100/json/").Return(nil, errors.New("connection error"))

	// Cria service com HTTP client mockado
	service := NewCEPServiceWithClient(mockClient)

	// Executa o teste
	result, err := service.GetLocationByCEP("01310-100")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching CEP data")

	// Verifica se o mock foi chamado corretamente
	mockClient.AssertExpectations(t)
} 