package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"gopher-mart/internal/utils"
	"net/http"
)

type balanceWithdrawHandler struct {
	usecase balanceWithdrawUsecase
}

func NewBalanceWithdrawHandler(usecase balanceWithdrawUsecase) *balanceWithdrawHandler {
	return &balanceWithdrawHandler{usecase: usecase}
}

type balanceWithdrawUsecase interface {
	WithdrawUserBonuses(ctx context.Context, user *users.User, wd *withdraws.Withdraw) error
	usecaseUsers.UserContextUsecase
}

type withdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

func (h *balanceWithdrawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Err(err).Send()
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Err(errors.ErrWrongContentType).Send()
		return
	}

	var req *withdrawRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	// convert to Withdraw
	wd := &withdraws.Withdraw{
		Sum:     req.Sum,
		OrderID: req.OrderID,
	}

	if !utils.LuhnValidator(req.OrderID) {
		fmt.Println()
		log.Warn().Str("id", req.OrderID)
		fmt.Println()
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.usecase.WithdrawUserBonuses(r.Context(), user, wd)
	if err != nil {
		log.Err(err).Send()
		switch err {
		case errors.ErrNotEnoughBonuses:
			w.WriteHeader(http.StatusPaymentRequired)
		case errors.ErrOrderWrongFormat:
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
