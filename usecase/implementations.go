package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/internal/auth"
	"github.com/SawitProRecruitment/UserService/internal/config"
	"github.com/SawitProRecruitment/UserService/internal/hashutil"
)

type AuthUsecaseImpl struct {
	userRepo UserRepository
	cfg      config.Config
}

// NewAuthUsecaseImpl instantiates AuthUsecaseImpl
func NewAuthUsecaseImpl(userRepo UserRepository, cfg config.Config) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// GetProfile fetches a user profile
func (u *AuthUsecaseImpl) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	return u.userRepo.Get(ctx, userID)
}

// Login logins a user
func (u *AuthUsecaseImpl) Login(ctx context.Context, phoneNumber, password string) (user *entity.User, accessToken *auth.AccessToken, err error) {
	user, err = u.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return
	}

	if !hashutil.IsPasswordMatched(password, user.HashedPassword) {
		return user, nil, entity.ErrInvalidPassword
	}

	accessToken, err = auth.GenerateAccessToken(user.ID, u.cfg.SecretKey)
	if err != nil {
		return
	}

	_ = u.userRepo.IncrementLoginCount(ctx, user.ID)

	return user, accessToken, nil
}

// Register registers a new user
func (u *AuthUsecaseImpl) Register(ctx context.Context, params CreateParams) (*entity.User, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashutil.HashSaltPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.Create(ctx, params, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile updates a user profile
func (u *AuthUsecaseImpl) UpdateProfile(ctx context.Context, params UpdateParams) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	_, err = u.userRepo.GetByPhoneNumber(ctx, params.PhoneNumber)
	if err == nil {
		return entity.ErrPhoneAlreadyExists
	}
	if err != entity.ErrUserNotFound {
		return err
	}

	return u.userRepo.Update(ctx, params)
}
