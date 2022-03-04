package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kolakanmi/grey_wallet/pkg/http/response"
	"github.com/Kolakanmi/grey_wallet/service"
)

type Handler struct {
	service service.ITransactionService
}

func New(service service.ITransactionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Credit(w http.ResponseWriter, r *http.Request) error {
	log.Println("Credit")
	var req service.TxnRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	result, err := h.service.Credit(r.Context(), req.Amount)
	if err != nil {
		return err
	}
	return response.OK("Success", result).ToJSON(w)
}

func (h *Handler) Debit(w http.ResponseWriter, r *http.Request) error {
	var req service.TxnRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	result, err := h.service.Debit(r.Context(), req.Amount)
	if err != nil {
		return err
	}
	return response.OK("Success", result).ToJSON(w)
}

func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) error {
	result, err := h.service.Balance(r.Context())
	if err != nil {
		return err
	}
	return response.OK("Success", result).ToJSON(w)
}

func (h *Handler) GetAllTransactions(w http.ResponseWriter, r *http.Request) error {
	result, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}
	return response.OK("Success", result).ToJSON(w)
}
