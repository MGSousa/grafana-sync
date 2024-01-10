package grafana

import (
	"net/http"

	"github.com/grafana-tools/sdk"
)

var httpClient = sdk.DefaultHTTPClient

type customHttpTransport struct {
	http.Transport

	customHeaders map[string]string
}

// RoundTrip implements a RoundTripper over HTTP to inject custom headers
// overrides the default one on HTTP.Transport
func (ct *customHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range ct.customHeaders {
		req.Header.Add(key, value)
	}

	return ct.Transport.RoundTrip(req)
}

// InitHttpClient
func InitHttpClient(customHeaders map[string]string) {
	httpClient.Transport = &customHttpTransport{
		customHeaders: customHeaders,
	}
}
