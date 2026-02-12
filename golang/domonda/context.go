package domonda

import (
	"context"
)

type ctxBaseURL struct{}

// WithBaseURL creates context.Context with baseURL param for making requests
func WithBaseURL(ctx context.Context, baseURL string) context.Context {
	return context.WithValue(ctx, ctxBaseURL{}, baseURL)
}

func baseURLFromCtx(ctx context.Context) string {
	baseURL, ok := ctx.Value(ctxBaseURL{}).(string)
	if !ok || baseURL == "" {
		return BaseURL
	}

	return baseURL
}
