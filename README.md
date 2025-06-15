# Weather CEP API

[![Deploy to Cloud Run](https://img.shields.io/badge/Deploy%20to-Cloud%20Run-blue)](https://cloud.google.com/run/docs?hl=pt-br)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

API em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em Celsius, Fahrenheit e Kelvin).

**🚀 Projeto desenvolvido como desafio do Lab 01: Deploy com Cloud Run**

## 📋 **Requisitos**

- Go 1.21+
- Chave da API [WeatherAPI](https://www.weatherapi.com/)
- Docker (opcional)

## 🚀 **Configuração Local**

### 1. Clone o repositório
```bash
git clone https://github.com/danielencestari/lab01.git
cd lab01
```

### 2. Configure as variáveis de ambiente
```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite o .env e adicione sua chave da WeatherAPI
WEATHER_API_KEY=sua_chave_aqui
```

**📝 Como obter a chave da WeatherAPI:**
1. Acesse [WeatherAPI.com](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Copie sua API key do dashboard
4. Cole no arquivo `.env`

### 3. Instale as dependências
```bash
go mod tidy
```

### 4. Execute a aplicação
```bash
# Opção 1: Comando direto
go run main.go

# Opção 2: Usando Make
make run
```

## 🧪 **Testando a API**

### Endpoints disponíveis:

- **Health Check**: `GET /health`
- **Temperatura por CEP**: `GET /temperature/{cep}`

### Exemplos de uso:

```bash
# Health check
curl http://localhost:8080/health

# Consulta CEP de São Paulo (Av. Paulista)
curl http://localhost:8080/temperature/01310100

# Consulta CEP do Rio de Janeiro (Copacabana)
curl http://localhost:8080/temperature/22070900

# Consulta CEP de Brasília (Asa Norte)
curl http://localhost:8080/temperature/70040010

# Resposta esperada (exemplo):
# {"temp_C":25.0,"temp_F":77.0,"temp_K":298.0}
```

## ☁️ **Deploy no Google Cloud Run**

### URL do Deploy
A API está disponível em: [https://weather-api-208690729789.us-central1.run.app/health](https://weather-api-208690729789.us-central1.run.app/health)

### Endpoints Disponíveis
- Health Check: `GET https://weather-api-208690729789.us-central1.run.app/health`
- Temperatura por CEP: `GET https://weather-api-208690729789.us-central1.run.app/temperature/{cep}`


### Exemplo de Resposta
```json
{
    "cidade": "São Paulo",
    "temp_C": 21.1,
    "temp_F": 69.98,
    "temp_K": 294.1
}
```



