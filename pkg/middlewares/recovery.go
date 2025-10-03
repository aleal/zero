package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

// Recovery middleware recovers from panics
func Recovery() Middleware {
	return func(next request.Handler) request.Handler {
		return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger := log.FromContext(rctx)
					logger.Error(rctx, "Panic recovered: %v", err)
					response.WriteError(w, http.StatusInternalServerError,
						fmt.Errorf("internal server error \n\n %v", err))
				}
			}()

			next(rctx, w, r)
		}
	}
}
