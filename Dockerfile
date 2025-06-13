# Build stage
FROM golang:1.21 as build
WORKDIR /app
COPY . .
COPY go.mod go.sum ./
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]
