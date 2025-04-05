package feign

import (
	"net/http"
	"sync"
	"time"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type TokenProvider interface {
	GetToken(baseURL string) string
}

type CachedTokenProvider struct {
	mu          sync.Mutex
	tokens      map[string]*Token
	fetchFunc   func(baseURL string) (string, time.Duration, error)
	refreshTime time.Duration // tempo antes da expiração para renovar
}

func NewCachedTokenProvider(fetch func(baseURL string) (string, time.Duration, error)) *CachedTokenProvider {
	return &CachedTokenProvider{
		tokens:      make(map[string]*Token),
		fetchFunc:   fetch,
		refreshTime: 10 * time.Second,
	}
}

func (p *CachedTokenProvider) GetToken(baseURL string) string {
	p.mu.Lock()
	defer p.mu.Unlock()

	token := p.tokens[baseURL]
	if token == nil || time.Now().After(token.ExpiresAt.Add(-p.refreshTime)) {
		newToken, expiresIn, err := p.fetchFunc(baseURL)
		if err != nil {
			// opcional: logar o erro e manter o token atual
			if token != nil {
				return token.Value
			}
			return ""
		}
		p.tokens[baseURL] = &Token{
			Value:     newToken,
			ExpiresAt: time.Now().Add(expiresIn),
		}
	}

	return p.tokens[baseURL].Value
}

type TokenTransport struct {
	Transport     http.RoundTripper
	TokenProvider TokenProvider
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	baseURL := getBaseURL(req.URL.String())
	token := t.TokenProvider.GetToken(baseURL)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return t.Transport.RoundTrip(req)
}

func getBaseURL(fullURL string) string {
	// Remove caminho e query para obter o host base
	u, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return ""
	}
	scheme := "http"
	if u.URL.Scheme != "" {
		scheme = u.URL.Scheme
	}
	return scheme + "://" + u.URL.Host
}
