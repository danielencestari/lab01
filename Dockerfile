# Build stage
FROM golang:1.21-alpine AS builder

# Instala certificados SSL
RUN apk --no-cache add ca-certificates

# Define diretório de trabalho
WORKDIR /app

# Copia arquivos de dependência
COPY go.mod go.sum ./

# Baixa dependências
RUN go mod download

# Copia código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Instala certificados SSL
RUN apk --no-cache add ca-certificates

# Cria usuário não-root
RUN adduser -D -s /bin/sh appuser

# Define diretório de trabalho
WORKDIR /root/

# Copia binário do estágio de build
COPY --from=builder /app/main .

# Muda para usuário não-root
USER appuser

# Expõe porta 8080
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 