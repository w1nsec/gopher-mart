package middlewares

import (
	"golang.org/x/net/context"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/users"
	"net/http"
)

type AuthMiddleware struct {
	usecase AuthUsecase
}

func NewAuthMidleware(usecase AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{usecase: usecase}
}

type AuthUsecase interface {
	Auth(ctx context.Context, user *users.User) error
	GenUserWithCookie(cookie *http.Cookie) *users.User
}

func (m *AuthMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(domain.CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user := m.usecase.GenUserWithCookie(cookie)
		err = m.usecase.Auth(r.Context(), user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
