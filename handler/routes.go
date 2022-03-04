package handler

import (
	"net/http"

	"github.com/Kolakanmi/grey_wallet/pkg/http/handler"
	"github.com/Kolakanmi/grey_wallet/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/credit",
			Method:  http.MethodPost,
			Handler: handler.CustomHandler(h.Credit),
		},
		{
			Path:    "/debit",
			Method:  http.MethodPost,
			Handler: handler.CustomHandler(h.Debit),
		},
		{
			Path:    "/balance",
			Method:  http.MethodGet,
			Handler: handler.CustomHandler(h.Balance),
		},
		{
			Path:    "/getAllTransactions",
			Method:  http.MethodGet,
			Handler: handler.CustomHandler(h.GetAllTransactions),
		},
	}
}
