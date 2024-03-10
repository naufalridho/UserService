package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/internal/auth"
	"github.com/SawitProRecruitment/UserService/internal/liberr"
	"github.com/SawitProRecruitment/UserService/usecase"
)

// Register is an HTTP handler for register endpoint
func (s *Server) Register(r echo.Context) error {
	var req generated.RegisterRequest

	err := json.NewDecoder(r.Request().Body).Decode(&req)
	if err != nil {
		return r.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request"})
	}

	params := usecase.CreateParams{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}

	user, err := s.AuthUsecase.Register(r.Request().Context(), params)
	if err != nil {
		return responseError(r, err)
	}

	resp := generated.RegisterResponse{
		UserId:      user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	return r.JSON(http.StatusOK, resp)
}

// GetProfile is an HTTP handler for get profile endpoint
func (s *Server) GetProfile(r echo.Context) error {
	accessToken := auth.AccessTokenFromHeader(r.Request().Header)

	claims, err := auth.ParseAccessToken(accessToken, s.Config.SecretKey)
	if err != nil {
		return responseError(r, err)
	}

	userID := claims.ClientID

	user, err := s.AuthUsecase.GetProfile(r.Request().Context(), userID)
	if err != nil {
		return responseError(r, err)
	}

	resp := generated.UserResponse{
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
	}

	return r.JSON(http.StatusOK, resp)
}

// UpdateProfile is an HTTP handler for update profile endpoint
func (s *Server) UpdateProfile(r echo.Context) error {
	accessToken := auth.AccessTokenFromHeader(r.Request().Header)

	claims, err := auth.ParseAccessToken(accessToken, s.Config.SecretKey)
	if err != nil {
		return responseError(r, err)
	}

	var req generated.UpdateProfileRequest
	err = json.NewDecoder(r.Request().Body).Decode(&req)
	if err != nil {
		return r.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request"})
	}

	params := usecase.UpdateParams{UserID: claims.ClientID}
	if req.PhoneNumber != nil {
		params.PhoneNumber = *req.PhoneNumber
	}
	if req.FullName != nil {
		params.FullName = *req.FullName
	}

	err = s.AuthUsecase.UpdateProfile(r.Request().Context(), params)
	if err != nil {
		return responseError(r, err)
	}

	return r.JSON(http.StatusOK, generated.GeneralResponse{Message: "Updated successfully"})
}

// Login is an HTTP handler for login endpoint
func (s *Server) Login(r echo.Context) error {
	var req generated.LoginRequest

	err := json.NewDecoder(r.Request().Body).Decode(&req)
	if err != nil {
		return r.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request"})
	}

	user, accessToken, err := s.AuthUsecase.Login(r.Request().Context(), req.PhoneNumber, req.Password)
	if err != nil {
		return responseError(r, err)
	}

	resp := generated.LoginResponse{
		UserId:      user.ID,
		AccessToken: accessToken.Value,
		Lifetime:    int(accessToken.Lifetime.Seconds()),
	}

	return r.JSON(http.StatusOK, resp)
}

func responseError(r echo.Context, err error) error {
	lerr, isLerr := err.(*liberr.Error)
	if !isLerr {
		return err
	}

	m := map[string]int{
		entity.ErrCodeUserNotFound:        http.StatusNotFound,
		entity.ErrCodePhoneAlreadyExists:  http.StatusConflict,
		entity.ErrCodeIneligiblePhone:     http.StatusBadRequest,
		entity.ErrCodePasswordLength:      http.StatusBadRequest,
		entity.ErrCodePasswordNumeric:     http.StatusBadRequest,
		entity.ErrCodePasswordCapital:     http.StatusBadRequest,
		entity.ErrCodePasswordSpecialChar: http.StatusBadRequest,
		entity.ErrCodeInvalidPassword:     http.StatusBadRequest,
		auth.ErrCodeInvalidJwt:            http.StatusForbidden,
	}

	statusCode, codeExists := m[lerr.Code]
	if !codeExists {
		statusCode = http.StatusInternalServerError
	}

	errResp := generated.ErrorResponse{
		Message: lerr.Error(),
	}

	if lerr.Field != "" {
		errResp.Field = &lerr.Field
	}

	return r.JSON(statusCode, errResp)
}
