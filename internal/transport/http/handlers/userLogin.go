package handlers

import (
	"context"
	"encoding/json"
	"gopher-mart/internal/domain/users"
	"net/http"
)

type LoginHandler struct {
	usecase loginUsecase
}

func NewLoginHandler(loginUsecase loginUsecase) *LoginHandler {
	return &LoginHandler{usecase: loginUsecase}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginUsecase interface {
	LoginUser(ctx context.Context, user *users.User) (cookie string, err error)
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := &users.User{
		Password: req.Password,
		Login:    req.Login,
	}

	cookie, err := h.usecase.LoginUser(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("set-cookie", cookie)
}