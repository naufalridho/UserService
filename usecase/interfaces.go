package usecase

import (
	"context"
	"github.com/SawitProRecruitment/UserService/entity"
)

type UserRepository interface {
	Create(ctx context.Context, params CreateParams, hashedPassword string) (user *entity.User, err error)
	Get(ctx context.Context, id string) (*entity.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)
	Update(ctx context.Context, params UpdateParams) error
	IncrementLoginCount(ctx context.Context, userID string) error
}
