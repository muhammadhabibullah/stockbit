package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	err = h.emitters[domain.DepositsTopic].EmitSync(key, value)
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

	key := fmt.Sprintf("%d", walletID)

	balanceView, err := h.viewers[domain.BalanceGroup].Get(key)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}
	if balanceView == nil {
		response := domain.GetDepositResponse{}
		data, _ := json.Marshal(response)
		_, _ = w.Write(data)
		return
	}

	balance, ok := balanceView.(*domain.Balance)
	if !ok {
		err = fmt.Errorf("unsupported balance type: %T", balanceView)
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	aboveThresholdView, err := h.viewers[domain.AboveThresholdGroup].Get(key)
	if err != nil {
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	aboveThreshold, ok := aboveThresholdView.(*domain.AboveThreshold)
	if !ok {
		err = fmt.Errorf("unsupported aboveThreshold type: %T", aboveThresholdView)
		http.Error(w, domain.NewHTTPError(err), http.StatusInternalServerError)
		return
	}

	const totalBalanceLimit float64 = 10000

	var (
		totalBalance float64
		timeLimit    = time.Now().Add(-2 * time.Minute)
	)

	for _, balanceHistory := range aboveThreshold.BalanceHistory {
		if balanceHistory.CreatedAt.Before(timeLimit) {
			break
		}

		totalBalance += balanceHistory.Amount
		if totalBalance > totalBalanceLimit {
			break
		}
	}

	response := domain.GetDepositResponse{
		Amount:         balance.Amount,
		AboveThreshold: totalBalance > totalBalanceLimit,
	}

	data, _ := json.Marshal(response)
	_, _ = w.Write(data)
}
