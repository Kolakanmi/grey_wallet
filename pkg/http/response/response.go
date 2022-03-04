package response

import (
	"encoding/json"
	"net/http"

	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
)

type (
	//RespBody - json response body format
	RespBody struct {
		Success bool        `json:"success"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}
	//Response - json response format
	Response struct {
		Body       *RespBody
		StatusCode int
	}
)

func newResponse(success bool, message string, data interface{}, status int, errorMessage string) *Response {
	return &Response{
		Body: &RespBody{
			Success: success,
			Message: message,
			Data:    data,
			Error:   errorMessage,
		},
		StatusCode: status,
	}
}

//OK returns response 200 status code
func OK(message string, data interface{}) *Response {
	return newResponse(true, message, data, http.StatusOK, "")
}

//Fail returns nil response body
func Fail(err error) *Response {
	errS := ""
	if err != nil {
		errS = err.Error()
	}
	appErr, ok := apperror.IsAppError(err)
	if ok {
		return newResponse(false, "", nil, appErr.Type(), errS)
	}
	return newResponse(false, "", nil, 500, errS)
}

//ToJSON converts response body to json
func (r *Response) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode)
	err := json.NewEncoder(w).Encode(r.Body)
	if err != nil {
		return err
	}
	return nil
}
