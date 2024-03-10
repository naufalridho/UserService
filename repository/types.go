// This file contains types that are used in the repository layer.
package repository

import (
	"database/sql"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/entity"
)

type userDto struct {
	ID             uuid.UUID    `db:"id"`
	FullName       string       `db:"full_name"`
	PhoneNumber    string       `db:"phone_number"`
	HashedPassword string       `db:"hashed_password"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
}

func (u userDto) toEntity() *entity.User {
	user := &entity.User{
		ID:             u.ID.String(),
		FullName:       u.FullName,
		PhoneNumber:    u.PhoneNumber,
		HashedPassword: u.HashedPassword,
		CreatedAt:      u.CreatedAt,
	}

	if u.UpdatedAt.Valid {
		user.UpdatedAt = &u.UpdatedAt.Time
	}

	return user
}

func userDtofromUpdateParams(params usecase.UpdateParams) (dto userDto, namedStmt string) {
	var userDto userDto
	var namedColumns []string

	userDto.ID = uuid.MustParse(params.UserID)

	if params.PhoneNumber != "" {
		userDto.PhoneNumber = params.PhoneNumber
		namedColumns = append(namedColumns, "phone_number = :phone_number")
	}

	if params.FullName != "" {
		userDto.FullName = params.FullName
		namedColumns = append(namedColumns, "full_name = :full_name")
	}

	return userDto, strings.Join(namedColumns, ", ")
}
