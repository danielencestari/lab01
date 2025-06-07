# Weather CEP API

API em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em Celsius, Fahrenheit e Kelvin).

## 📋 **Requisitos**

- Go 1.21+
- Chave da API [WeatherAPI](https://www.weatherapi.com/)
- Docker (opcional)

## 🚀 **Configuração Local**

### 1. Clone o repositório
```bash
git clone <seu-repositorio>
cd desafio_deploy_com_cloud_run
```

### 2. Configure as variáveis de ambiente
```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite o .env e adicione sua chave da WeatherAPI
WEATHER_API_KEY=sua_chave_aqui
```

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

# Consulta CEP de São Paulo
curl http://localhost:8080/temperature/01310100

# Resposta esperada:
# {"temp_C":25.0,"temp_F":77.0,"temp_K":298.0}
```

### Códigos de resposta:

- **200**: Sucesso
- **404**: CEP não encontrado
- **422**: CEP inválido
- **500**: Erro interno (verifique se a API key está configurada)

## 🧪 **Executando Testes**

```bash
# Todos os testes
make test

# Apenas testes unitários
make test-unit

# Testes com cobertura
make test-coverage

# Testar endpoints (aplicação deve estar rodando)
make test-endpoints
```

## 🐳 **Docker**

### Executar com Docker:
```bash
# Build da imagem
make docker-build

# Executar via Docker
make docker-run

# Ou com docker-compose
make docker-compose-up
```

## ☁️ **Deploy no Google Cloud Run**

### 1. **Configure as variáveis de ambiente no Cloud Run**

**⚠️ IMPORTANTE: Nunca coloque secrets no código ou dockerfile!**

No Google Cloud Console:
1. Vá para [Cloud Run](https://cloud.google.com/run/docs?hl=pt-br)
2. Deploy sua aplicação
3. Configure as variáveis de ambiente:
   - `WEATHER_API_KEY`: sua_chave_da_weatherapi
   - `GIN_MODE`: release

### 2. **Usando gcloud CLI:**

```bash
# Fazer build e deploy
gcloud run deploy weather-cep-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --set-env-vars="WEATHER_API_KEY=sua_chave_aqui,GIN_MODE=release" \
  --allow-unauthenticated
```

### 3. **Deploy automático via GitHub Actions:**

Veja o arquivo `.github/workflows/deploy.yml` para configuração de CI/CD.

## 📁 **Estrutura do Projeto**

```
├── handlers/           # HTTP handlers
├── services/          # Lógica de negócio
├── models/            # Estruturas de dados
├── utils/             # Funções utilitárias
├── main.go            # Ponto de entrada
├── Dockerfile         # Imagem Docker
├── docker-compose.yml # Configuração local
├── Makefile          # Comandos úteis
├── .env.example      # Template de configuração
└── README.md         # Este arquivo
```

## 🔒 **Segurança**

- ✅ Arquivo `.env` está no `.gitignore`
- ✅ Secrets são carregados via variáveis de ambiente
- ✅ Validação de entrada (CEP)
- ✅ Tratamento de erros apropriado

## 🚨 **Importante para Deploy**

1. **Nunca** commite arquivos `.env` ou chaves de API
2. Configure secrets via **variáveis de ambiente** no Cloud Run
3. Use **GIN_MODE=release** em produção
4. Configure **health checks** apropriados

## 📚 **APIs Utilizadas**

- [ViaCEP](https://viacep.com.br/) - Consulta de CEPs brasileiros
- [WeatherAPI](https://www.weatherapi.com/) - Dados meteorológicos

## 🛠 **Comandos Make Disponíveis**

```bash
make help          # Ver todos os comandos
make run           # Executar aplicação
make test          # Executar testes
make build         # Compilar aplicação
make docker-build  # Build Docker
make docker-run    # Executar via Docker
make clean         # Limpar arquivos de build
```

## 📝 **Licença**

Este projeto é licenciado sob a licença MIT. 