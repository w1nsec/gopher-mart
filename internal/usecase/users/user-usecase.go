package users

import (
	"context"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/repository"
	"net/http"
)

type UserUsecase interface {
	LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)
	RegisterUser(ctx context.Context, user *users.User) error

	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
}

type Usecase struct {
	storage repository.UserRepoInf
}

func NewUsecase(storage repository.UserRepoInf) *Usecase {
	return &Usecase{storage: storage}
}

func (u *Usecase) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	anyType := ctx.Value(domain.UserContextKey)
	user, ok := anyType.(*users.User)
	if !ok {
		return nil, errors.ErrUserNotFoundInContext
	}
	return user, nil
}

func (u *Usecase) CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error) {
	return u.storage.CheckBalance(ctx, user)
}

func (u *Usecase) LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	return u.storage.LoginUser(ctx, user)
}

func (u *Usecase) RegisterUser(ctx context.Context, user *users.User) error {
	return u.storage.RegisterUser(ctx, user)
}

func (u *Usecase) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	return u.storage.GetUserWithdrawals(ctx, user)

}
