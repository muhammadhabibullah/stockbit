package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"stockbit/domain"
)

func (h *httpHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getDeposit(w, r)
	case http.MethodPost:
		h.postDeposit(w, r)
	default:
		http.Error(w, domain.NewHTTPError(domain.ErrInvalidRequestMethod), http.StatusMethodNotAllowed)
		return
	}
}

func (h *httpHandler) postDeposit(w http.ResponseWriter, r *http.Request) {
	var req domain.DepositRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusBadRequest)
		return
	}

	if req.WalletID == 0 || req.Amount == 0 {
		err = fmt.Errorf("empty amount/wallet ID")
		http.Error(w, domain.NewHTTPError(err), http.StatusBadRequest)
		return
	}

	err = h.userUseCase.Deposit(r.Context(), req)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(&req, "", "    ")
	_, _ = fmt.Fprintf(w, string(response))
}

func (h *httpHandler) getDeposit(w http.ResponseWriter, r *http.Request) {
	walletIDQuery := r.URL.Query().Get("id")
	walletID, _ := strconv.ParseInt(walletIDQuery, 10, 64)
	if walletID == 0 {
		http.Error(w, domain.NewHTTPError(domain.ErrInvalidID), http.StatusBadRequest)
		return
	}

	response, err := h.userUseCase.GetDeposit(r.Context(), walletID)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(response)
	_, _ = w.Write(data)
}
