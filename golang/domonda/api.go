package domonda

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const (
	baseURL = "https://domonda.app/api/public"
)

func postJSON(ctx context.Context, apiKey, endpoint string, payload any) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", baseURL+endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(request)
}
