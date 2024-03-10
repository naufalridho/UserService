package entity

import "github.com/SawitProRecruitment/UserService/internal/liberr"

const (
	ErrCodeUserNotFound = "USER_SERVICE-USER-NOT_FOUND"

	ErrCodePhoneAlreadyExists = "USER_SERVICE-PHONE-ALREADY_EXISTS"
	ErrCodeIneligiblePhone    = "USER_SERVICE-PHONE-INELIGIBLE"

	ErrCodePasswordLength      = "USER_SERVICE-PASSWORD-LENGTH"
	ErrCodePasswordNumeric     = "USER_SERVICE-PASSWORD-NUMERIC"
	ErrCodePasswordCapital     = "USER_SERVICE-PASSWORD-CAPITAL"
	ErrCodePasswordSpecialChar = "USER_SERVICE-PASSWORD-SPECIAL_CHAR"
	ErrCodeInvalidPassword     = "USER_SERVICE-PASSWORD-INVALID"
)

var (
	ErrUserNotFound = liberr.NewError("User not found", ErrCodeUserNotFound, "id")

	ErrPhoneAlreadyExists = liberr.NewError("Phone number already exists", ErrCodePhoneAlreadyExists, "phone_number")
	ErrIneligiblePhone    = liberr.NewError("Phone number should contain +62 and has length between 10 and 13 chars", ErrCodeIneligiblePhone, "phone_number")

	ErrPasswordLength      = liberr.NewError("Password is too short or too long", ErrCodePasswordLength, "password")
	ErrPasswordNumeric     = liberr.NewError("Password must contain at least 1 numeric character", ErrCodePasswordNumeric, "password")
	ErrPasswordCapital     = liberr.NewError("Password must contain at least 1 capital character", ErrCodePasswordCapital, "password")
	ErrPasswordSpecialChar = liberr.NewError("Password must contain at least 1 special character", ErrCodePasswordSpecialChar, "password")
	ErrInvalidPassword     = liberr.NewError("Invalid password", ErrCodeInvalidPassword, "password")
)
