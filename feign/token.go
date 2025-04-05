package feign

import (
	"net/http"
)

type TokenProvider interface {
	GetToken() string
}

type TokenTransport struct {
	Transport     http.RoundTripper
	TokenProvider TokenProvider
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}
	token := t.TokenProvider.GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return t.Transport.RoundTrip(req)
}

