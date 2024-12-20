/*
Updating the master-data of a client-company can be done in bulk using the REST API
with JSON data. In general those API endpoints follow an `upsert` logic,
meaning that if data records can be identified by an ID or name,
they will be updated with the provided data,
else new records will be inserted.

For those endpoints we also provide a Go SDK with this package.

Using this Go package to make API requests has the benefit of basic
client side validation of the data before sending the request
and the client package will always be kept up to date
with the API server implementation.

The server is also implemented with Go and uses the standard library
JSON parser which is documented here: https://pkg.go.dev/encoding/json

The struct types in in this package define the JSON API
where every struct is mapped to a JSON object
with the names of the struct fields used as JSON object value names.
The Go JSON parser detects object values names case insensitive
and object values that are left out in the JSON will have
their types' Go zero value after parsing.
*/
package domonda

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	// BaseURL is the base URL for all API endpoints
	BaseURL = "https://domonda.app/api/public"

	// SourceTestEndpointNOP can be passed as magical value for
	// the source argument of the post functions
	// to test the endpoints without side effects.
	SourceTestEndpointNOP = "TestEndpointNOP"
)

func postJSON(ctx context.Context, apiKey, endpoint string, vals url.Values, payload any) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := BaseURL + endpoint
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
