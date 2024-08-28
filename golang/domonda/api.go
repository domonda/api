package domonda

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://domonda.app/api/public"
)

func postJSON(ctx context.Context, apiKey, endpoint string, vals url.Values, payload any) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := baseURL + endpoint
	if len(vals) > 0 {
		url += "?" + vals.Encode()
	}
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(request)
}
