package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
	"github.com/Kolakanmi/grey_wallet/pkg/http/response"
)

//Recover - middleware
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", err)
				}
				log.Printf("recover: %v", rec)
				stack := make([]byte, 4<<10)
				length := runtime.Stack(stack, false)
				log.Printf("stack: %s, err: %v", string(stack[:length]), err)
				response.Fail(apperror.CouldNotCompleteRequest()).ToJSON(w)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
