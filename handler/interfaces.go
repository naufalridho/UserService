package handler

import (
	"context"
	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/internal/auth"
	"github.com/SawitProRecruitment/UserService/usecase"
)

type AuthUsecase interface {
	Register(ctx context.Context, params usecase.CreateParams) (*entity.User, error)
	Login(ctx context.Context, userID, password string) (user *entity.User, accessToken *auth.AccessToken, err error)
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, params usecase.UpdateParams) error
}
