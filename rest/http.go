package rest

import (
	"context"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HttpClientDoer implements HttpDoer
type HttpClientDoer struct {
	HttpClient HttpClient
}

// Do overrides Do method of the default goa client Doer. It's needed for mocking http clients in tests.
func (d *HttpClientDoer) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return d.HttpClient.Do(req)
}
