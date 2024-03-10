package fixtures

import (
	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/usecase"
	"time"
)

const (
	userID      = "801868c6-a641-46cc-b5bf-2a8956ef14af"
	fullName    = "John Smith"
	phoneNumber = "+628123456789"
	secretKey   = "QcRPpsGwuHNAoWvOrWmM"
)

func User() *entity.User {
	return &entity.User{
		ID:             userID,
		FullName:       fullName,
		PhoneNumber:    phoneNumber,
		HashedPassword: "$2a$04$0iDToJ7e0u2oagELyk3s.uRGZgyz1rMaKlubO6pvHSMeXD5brNj.m",
		CreatedAt:      time.Now(),
	}
}

func CreateParams() usecase.CreateParams {
	return usecase.CreateParams{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Password:    "5awitPro!",
	}
}

func UpdateParams() usecase.UpdateParams {
	return usecase.UpdateParams{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: "+6281987654321",
	}
}

func SecretKey() string {
	return secretKey
}
