package handler_test

import (
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/internal/auth"
	"github.com/SawitProRecruitment/UserService/internal/test/fixtures"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/SawitProRecruitment/UserService/internal/config"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/handler/mocks"
)

func TestServer_RegisterSuccess(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"password\": \"5awitPro!\", \"phone_number\": \"+628123456789\"}"

	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().Register(r.Request().Context(), fixtures.CreateParams()).Return(fixtures.User(), nil)

	h := handler.NewServer(au, config.Config{SecretKey: "QcRPpsGwuHNAoWvOrWmM"})

	// Assertions
	if assert.NoError(t, h.Register(r)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestServer_RegisterInvalidJsonError(t *testing.T) {
	e := echo.New()

	payload := "{xxxx}"

	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)
	au := mocks.NewMockAuthUsecase(ctrl)

	h := handler.NewServer(au, config.Config{SecretKey: "QcRPpsGwuHNAoWvOrWmM"})

	// Assertions
	if assert.NoError(t, h.Register(r)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestServer_RegisterUsecaseError(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"password\": \"5awitPro!\", \"phone_number\": \"+628123456789\"}"

	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().Register(r.Request().Context(), fixtures.CreateParams()).Return(nil, errors.New("some error"))

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	assert.Error(t, h.Register(r))
}

func TestServer_GetProfileSuccess(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"phone_number\": \"+628123456789\"}"

	user := fixtures.User()
	token, _ := auth.GenerateAccessToken(user.ID, fixtures.SecretKey())

	req := httptest.NewRequest(http.MethodGet, "/v1/users", strings.NewReader(payload))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Value))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().GetProfile(r.Request().Context(), user.ID).Return(fixtures.User(), nil)

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	if assert.NoError(t, h.GetProfile(r)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestServer_GetProfileParseTokenError(t *testing.T) {
	e := echo.New()

	token := "xxxx"

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	if assert.NoError(t, h.UpdateProfile(r)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}
}

func TestServer_GetProfileUsecaseError(t *testing.T) {
	e := echo.New()

	user := fixtures.User()
	token, _ := auth.GenerateAccessToken(user.ID, fixtures.SecretKey())

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Value))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().GetProfile(r.Request().Context(), user.ID).Return(nil, errors.New("some error"))

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	assert.Error(t, h.GetProfile(r))
}

func TestServer_UpdateProfileSuccess(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"phone_number\": \"+6281987654321\"}"

	user := fixtures.User()
	token, _ := auth.GenerateAccessToken(user.ID, fixtures.SecretKey())

	req := httptest.NewRequest(http.MethodPatch, "/v1/users", strings.NewReader(payload))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Value))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().UpdateProfile(r.Request().Context(), fixtures.UpdateParams()).Return(nil)

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	if assert.NoError(t, h.UpdateProfile(r)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestServer_UpdateProfileInvalidJsonError(t *testing.T) {
	e := echo.New()

	payload := "xxxx"

	user := fixtures.User()
	token, _ := auth.GenerateAccessToken(user.ID, fixtures.SecretKey())

	req := httptest.NewRequest(http.MethodPatch, "/v1/users", strings.NewReader(payload))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Value))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	if assert.NoError(t, h.UpdateProfile(r)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestServer_UpdateProfileParseTokenError(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"phone_number\": \"+628123456789\"}"

	token := "xxxx"

	req := httptest.NewRequest(http.MethodPatch, "/v1/users", strings.NewReader(payload))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	if assert.NoError(t, h.UpdateProfile(r)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}
}

func TestServer_UpdateProfileUsecaseError(t *testing.T) {
	e := echo.New()

	payload := "{\"full_name\": \"John Smith\", \"phone_number\": \"+6281987654321\"}"

	user := fixtures.User()
	token, _ := auth.GenerateAccessToken(user.ID, fixtures.SecretKey())

	req := httptest.NewRequest(http.MethodPatch, "/v1/users", strings.NewReader(payload))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Value))

	rec := httptest.NewRecorder()
	r := e.NewContext(req, rec)

	ctrl := gomock.NewController(t)

	au := mocks.NewMockAuthUsecase(ctrl)
	au.EXPECT().UpdateProfile(r.Request().Context(), fixtures.UpdateParams()).Return(errors.New("some error"))

	h := handler.NewServer(au, config.Config{SecretKey: fixtures.SecretKey()})

	// Assertions
	assert.Error(t, h.UpdateProfile(r))
}
