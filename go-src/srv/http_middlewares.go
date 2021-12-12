package srv

import (
	"fmt"
	"net/http"
)

func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				switch recValue := rec.(type) {
				case error:
					writeError(w, recValue)
				case string:
					writeError(w, fmt.Errorf(recValue))
				default:
					writeError(w, fmt.Errorf("internal error"))
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
