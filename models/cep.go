package models

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        string `json:"erro,omitempty"`
}

// LocationInfo representa as informações de localização extraídas do CEP
type LocationInfo struct {
	City  string `json:"city"`
	State string `json:"state"`
	CEP   string `json:"cep"`
} 