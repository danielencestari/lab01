package services

import "net/http"
 
// HTTPClientInterface define o contrato para HTTP clients
type HTTPClientInterface interface {
	Get(url string) (*http.Response, error)
} 