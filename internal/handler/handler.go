package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	info "github.com/Dimadetected/my-bank-service/interface"
	"github.com/Dimadetected/my-bank-service/internal/service"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

type SumRequest struct {
	Sum float64 `json:"sum"`
}
type CurrencyRequest struct {
	Currency string `json:"currency"`
}

func (h *Handler) AddFunds(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}
	fmt.Println(string(body))
	var b SumRequest
	if err := json.Unmarshal(body, &b); err != nil {
		WriteResp(rw, http.StatusInternalServerError, err.Error())
	}
	h.s.Account.AddFunds(b.Sum)
	WriteResp(rw, http.StatusOK, "сумма успешно добавлена")
}
func (h *Handler) Withdraw(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}

	var b SumRequest
	if err := json.Unmarshal(body, &b); err != nil {
		WriteResp(rw, http.StatusInternalServerError, err.Error())
	}

	if err := h.s.Account.Withdraw(b.Sum); err != nil {
		WriteResp(rw, http.StatusInternalServerError, err.Error())
	}
	WriteResp(rw, http.StatusOK, "списание произошло успешно")

}
func (h *Handler) GetCurrency(rw http.ResponseWriter, r *http.Request) {
	curr := h.s.Account.GetCurrency()
	WriteResp(rw, http.StatusOK, curr)
}
func (h *Handler) SumProfit(rw http.ResponseWriter, r *http.Request) {
	h.s.Account.SumProfit()
	WriteResp(rw, http.StatusOK, "начисления успешно посчитаны")

}
func (h *Handler) GetAccountCurrencyRate(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}

	var b CurrencyRequest
	if err := json.Unmarshal(body, &b); err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}

	current := h.s.Account.GetAccountCurrencyRate(info.Currency(b.Currency))
	WriteResp(rw, http.StatusOK, fmt.Sprintf("%f %s", current, b.Currency))

}
func (h *Handler) GetBalance(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}

	var b CurrencyRequest
	if err := json.Unmarshal(body, &b); err != nil {
		WriteResp(rw, http.StatusBadRequest, err.Error())
	}

	balance := h.s.Account.GetBalance(info.Currency(b.Currency))
	WriteResp(rw, http.StatusOK, fmt.Sprintf("%f %s", balance, "SBP"))
}
func WriteResp(rw http.ResponseWriter, code int, errStr interface{}) {
	rw.WriteHeader(code)
	m := make(map[string]interface{})
	if code == http.StatusOK {
		m["answer"] = errStr
	} else {
		m["error"] = errStr
	}
	resp, err := json.Marshal(m)
	if err != nil {
		resp = []byte("err")
	}

	rw.Write(resp)
}
