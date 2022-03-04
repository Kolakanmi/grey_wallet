package handler

import (
	"log"
	"net/http"

	"github.com/Kolakanmi/grey_wallet/pkg/http/response"
)

//CustomHandler - implements the Handler interface and logs error/ returns error
type CustomHandler func(http.ResponseWriter, *http.Request) error

func (c CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c(w, r)
	if err == nil {
		return
	}

	err = respondWithError(err, w)
	if err != nil {
		log.Printf("error decoding json: %v", err)
	}

}

func respondWithError(err error, w http.ResponseWriter) error {
	return response.Fail(err).ToJSON(w)
}
