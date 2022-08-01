package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/proto"
	"stockbit/domain"
	"stockbit/domain/proto/pb"
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

	depositEmit := pb.Deposit{
		WalletId: req.WalletID,
		Amount:   float32(req.Amount),
	}
	msg, err := proto.Marshal(&depositEmit)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	key := fmt.Sprint(req.WalletID)
	value := string(msg)

	err = h.emitter.EmitSync(key, value)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(&depositEmit, "", "    ")
	_, _ = fmt.Fprintf(w, string(response))
}

func (h *httpHandler) getDeposit(w http.ResponseWriter, r *http.Request) {
	walletIDQuery := r.URL.Query().Get("id")
	walletID, _ := strconv.ParseInt(walletIDQuery, 10, 64)
	if walletID == 0 {
		http.Error(w, domain.NewHTTPError(domain.ErrInvalidID), http.StatusBadRequest)
		return
	}

	value, err := h.view.Get(fmt.Sprintf("%d", walletID))
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(value)
	_, _ = w.Write(data)
}
