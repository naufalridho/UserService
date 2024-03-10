package usecase

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/entity"
	"regexp"
)

const (
	minFullNameLen = 3
	maxFullNameLen = 60

	minPhoneNumberLen = 10
	maxPhoneNumberLen = 13

	minPasswordLen = 6
	maxPasswordLen = 64
)

// CreateParams holds user creation parameters
type CreateParams struct {
	FullName    string
	PhoneNumber string
	Password    string
}

// Validate validates password and phone number in the create parameters
func (p CreateParams) Validate() error {
	err := validatePassword(p.Password)
	if err != nil {
		return err
	}

	return validatePhone(p.PhoneNumber)
}

// UpdateParams holds user update parameters
type UpdateParams struct {
	UserID      string
	FullName    string
	PhoneNumber string
}

// Validate validates phone number and full name in the update parameters
func (u UpdateParams) Validate() error {
	if err := validatePhone(u.PhoneNumber); u.PhoneNumber != "" && err != nil {
		return err
	}

	if err := validateFullName(u.FullName); u.FullName != "" && err != nil {
		return err
	}

	return nil
}

func validatePhone(phoneNumber string) error {
	// This checks a password whether a phone number has +62 prefix
	// and phone number min and max length
	regex, _ := regexp.Compile(fmt.Sprintf(`^\+62\d{%d,%d}$`, minPhoneNumberLen, maxPhoneNumberLen))

	if !regex.MatchString(phoneNumber) {
		return entity.ErrIneligiblePhone
	}

	return nil
}

func validateFullName(fullName string) error {
	if len(fullName) < minFullNameLen || len(fullName) > maxFullNameLen {
		return entity.ErrPhoneAlreadyExists
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return entity.ErrPasswordLength
	}

	// This checks a password whether it is containing at least 1 capital character
	regex1, _ := regexp.Compile(`.*[A-Z].*`)

	if !regex1.MatchString(password) {
		return entity.ErrPasswordCapital
	}

	// This checks a password whether it is containing at least 1 numeric character
	regex2, _ := regexp.Compile(`.*\d.*`)

	if !regex2.MatchString(password) {
		return entity.ErrPasswordNumeric
	}

	// This checks a password whether it is containing at least 1 special character
	regex3, _ := regexp.Compile(`.*[^A-Za-z\d].*`)

	if !regex3.MatchString(password) {
		return entity.ErrPasswordSpecialChar
	}

	return nil
}
