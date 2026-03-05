package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HumaInternalAuth rejects requests missing the correct X-Internal-Secret header.
func HumaInternalAuth(secret string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if ctx.Header("X-Internal-Secret") != secret {
			ctx.SetStatus(http.StatusForbidden)
			_, _ = ctx.BodyWriter().Write([]byte(`{"title":"Forbidden","status":403,"detail":"invalid internal secret"}`))
			return
		}
		next(ctx)
	}
}

// HumaRequireUser returns 401 when X-User-ID is absent from context.
func HumaRequireUser(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		userID := GetUserID(ctx.Context())
		if userID == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "authentication required")
			return
		}
		next(ctx)
	}
}
