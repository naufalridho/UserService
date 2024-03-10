package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/internal/config"
	"github.com/SawitProRecruitment/UserService/internal/test/fixtures"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/SawitProRecruitment/UserService/usecase/mocks"
)

type authUsecaseTestSuite struct {
	subject  *usecase.AuthUsecaseImpl
	userRepo *mocks.MockUserRepository
	config   config.Config
}

func TestNewAuthUsecaseImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := usecase.NewAuthUsecaseImpl(userRepo, fixtures.Config())

	assert.NotNil(t, uc)
}

func TestAuthUsecaseImpl_GetProfile(t *testing.T) {
	for _, test := range []struct {
		name string
		fn   func(*testing.T, *authUsecaseTestSuite)
	}{
		{
			name: "success case",
			fn:   testGetProfileSuccess,
		},
		{
			name: "error case - get user error",
			fn:   testGetProfileGetError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)

			suite := &authUsecaseTestSuite{
				subject:  usecase.NewAuthUsecaseImpl(mockUserRepo, fixtures.Config()),
				userRepo: mockUserRepo,
			}

			test.fn(t, suite)
		})
	}
}

func testGetProfileSuccess(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	userID := uuid.New().String()

	user := fixtures.User()

	suite.userRepo.EXPECT().Get(ctx, userID).Return(user, nil)
	got, err := suite.subject.GetProfile(ctx, userID)

	assert.Equal(t, user, got)
	assert.Nil(t, err)
}

func testGetProfileGetError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	userID := uuid.New().String()

	suite.userRepo.EXPECT().Get(ctx, userID).Return(nil, errors.New("some error"))

	_, err := suite.subject.GetProfile(ctx, userID)
	assert.NotNil(t, err)
}

func TestAuthUsecaseImpl_Login(t *testing.T) {
	for _, test := range []struct {
		name string
		fn   func(*testing.T, *authUsecaseTestSuite)
	}{
		{
			name: "success case",
			fn:   testLoginSuccess,
		},
		{
			name: "error case - get user error",
			fn:   testLoginGetUserError,
		},
		{
			name: "error case - invalid password ",
			fn:   testLoginInvalidPasswordError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)

			suite := &authUsecaseTestSuite{
				subject:  usecase.NewAuthUsecaseImpl(mockUserRepo, fixtures.Config()),
				userRepo: mockUserRepo,
			}

			test.fn(t, suite)
		})
	}
}

func testLoginSuccess(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	password := "5awitPro!"
	user := fixtures.User()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, user.PhoneNumber).Return(user, nil)
	suite.userRepo.EXPECT().IncrementLoginCount(ctx, user.ID).Return(nil)

	_, got, err := suite.subject.Login(ctx, user.PhoneNumber, password)

	assert.NotEmpty(t, got)
	assert.Nil(t, err)
}

func testLoginGetUserError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	password := "5awitPro!"
	user := fixtures.User()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, user.PhoneNumber).Return(nil, errors.New("some error"))
	_, _, err := suite.subject.Login(ctx, user.PhoneNumber, password)

	assert.NotNil(t, err)
}

func testLoginInvalidPasswordError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	password := "somewrongpassword"
	user := fixtures.User()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, user.PhoneNumber).Return(user, nil)
	_, _, err := suite.subject.Login(ctx, user.PhoneNumber, password)

	assert.Equal(t, entity.ErrInvalidPassword, err)
}

func TestAuthUsecaseImpl_UpdateProfile(t *testing.T) {
	for _, test := range []struct {
		name string
		fn   func(*testing.T, *authUsecaseTestSuite)
	}{
		{
			name: "success case",
			fn:   testUpdateProfileSuccess,
		},
		{
			name: "error case - phone already exists error",
			fn:   testUpdateProfileAlreadyExistsError,
		},
		{
			name: "error case - full name length",
			fn:   testUpdateProfileFullNameLengthError,
		},
		{
			name: "error case - get user error",
			fn:   testUpdateProfileGetUserError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)

			suite := &authUsecaseTestSuite{
				subject:  usecase.NewAuthUsecaseImpl(mockUserRepo, fixtures.Config()),
				userRepo: mockUserRepo,
			}

			test.fn(t, suite)
		})
	}
}

func testUpdateProfileSuccess(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	params := fixtures.UpdateParams()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, params.PhoneNumber).Return(nil, entity.ErrUserNotFound)
	suite.userRepo.EXPECT().Update(ctx, params).Return(nil)

	err := suite.subject.UpdateProfile(ctx, params)

	assert.Nil(t, err)
}

func testUpdateProfileAlreadyExistsError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	params := fixtures.UpdateParams()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, params.PhoneNumber).Return(fixtures.User(), nil)

	err := suite.subject.UpdateProfile(ctx, params)

	assert.Equal(t, entity.ErrPhoneAlreadyExists, err)
}

func testUpdateProfileFullNameLengthError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	params := fixtures.UpdateParams()

	params.FullName = "a"

	err := suite.subject.UpdateProfile(ctx, params)

	assert.Equal(t, entity.ErrPhoneAlreadyExists, err)
}

func testUpdateProfileGetUserError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()
	params := fixtures.UpdateParams()

	suite.userRepo.EXPECT().GetByPhoneNumber(ctx, params.PhoneNumber).Return(nil, errors.New("some error"))

	err := suite.subject.UpdateProfile(ctx, params)

	assert.NotNil(t, err)
}

func TestAuthUsecaseImpl_Register(t *testing.T) {
	for _, test := range []struct {
		name string
		fn   func(*testing.T, *authUsecaseTestSuite)
	}{
		{
			name: "success case",
			fn:   testRegisterSuccess,
		},
		{
			name: "error case - password does not contain capital character",
			fn:   testRegisterPasswordNotContainCapitalError,
		},
		{
			name: "error case - password does not contain numeric character",
			fn:   testRegisterPasswordNotContainNumericError,
		},
		{
			name: "error case - password does not contain special character",
			fn:   testRegisterPasswordNotContainSpecialCharError,
		},
		{
			name: "error case - password is too long or too short",
			fn:   testRegisterPasswordLengthError,
		},
		{
			name: "error case - phone number is not eligible",
			fn:   testRegisterIneligiblePhoneError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)

			suite := &authUsecaseTestSuite{
				subject:  usecase.NewAuthUsecaseImpl(mockUserRepo, fixtures.Config()),
				userRepo: mockUserRepo,
			}

			test.fn(t, suite)
		})
	}
}

func testRegisterSuccess(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()
	user := fixtures.User()

	suite.userRepo.EXPECT().Create(ctx, params, gomock.Not("")).Return(user, nil)
	got, err := suite.subject.Register(ctx, params)

	assert.NotEmpty(t, got)
	assert.Nil(t, err)
}

func testRegisterPasswordNotContainCapitalError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()

	params.Password = "5awitpro!"

	_, err := suite.subject.Register(ctx, params)

	assert.Equal(t, entity.ErrPasswordCapital, err)
}

func testRegisterPasswordNotContainNumericError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()

	params.Password = "SawitPro!"

	_, err := suite.subject.Register(ctx, params)

	assert.Equal(t, entity.ErrPasswordNumeric, err)
}

func testRegisterPasswordNotContainSpecialCharError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()

	params.Password = "5awitPro"

	_, err := suite.subject.Register(ctx, params)

	assert.Equal(t, entity.ErrPasswordSpecialChar, err)
}

func testRegisterPasswordLengthError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()

	params.Password = "S0!"

	_, err := suite.subject.Register(ctx, params)

	assert.Equal(t, entity.ErrPasswordLength, err)
}

func testRegisterIneligiblePhoneError(t *testing.T, suite *authUsecaseTestSuite) {
	ctx := context.Background()

	params := fixtures.CreateParams()

	params.PhoneNumber = "123456789"

	_, err := suite.subject.Register(ctx, params)

	assert.Equal(t, entity.ErrIneligiblePhone, err)
}
