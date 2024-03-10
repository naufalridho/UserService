package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/usecase"
)

const (
	ErrCodeConflict = "23505"
)

// UserRepository holds user repository
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, params usecase.CreateParams, hashedPassword string) (user *entity.User, err error) {
	newId := uuid.New()
	query := "INSERT INTO users (id, full_name, phone_number, hashed_password) VALUES ($1,$2,$3,$4)"

	_, err = r.db.ExecContext(ctx, query, newId, params.FullName, params.PhoneNumber, hashedPassword)
	if err == nil {
		return &entity.User{
			ID:             newId.String(),
			FullName:       params.FullName,
			PhoneNumber:    params.PhoneNumber,
			HashedPassword: hashedPassword,
		}, nil
	}

	if pqErr, isPqErr := err.(*pq.Error); isPqErr && pqErr.Code == ErrCodeConflict {
		return nil, entity.ErrPhoneAlreadyExists
	}

	log.Error("[UserRepository][Create] an error occurred: ", err)
	return nil, err
}

// Get fetches a user by his user ID
func (r *UserRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	var userDto userDto

	query := "SELECT id, full_name, phone_number, hashed_password, created_at, updated_at FROM users WHERE id = $1"
	row := r.db.QueryRowxContext(ctx, query, id)

	err := row.StructScan(&userDto)
	if err != nil {
		log.Error("[UserRepository][Get] an error occurred: ", err)
		return nil, err
	}

	return userDto.toEntity(), nil
}

// GetByPhoneNumber fetches a user by his phone number
func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	var userDto userDto

	query := "SELECT id, full_name, phone_number, hashed_password, created_at, updated_at FROM users WHERE phone_number = $1"
	row := r.db.QueryRowxContext(ctx, query, phoneNumber)

	err := row.StructScan(&userDto)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrUserNotFound
	}

	if err != nil {
		log.Error("[UserRepository][GetByPhoneNumber] an error occurred: ", err)
		return nil, err
	}

	return userDto.toEntity(), nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, params usecase.UpdateParams) error {
	userDto, namedStmt := userDtofromUpdateParams(params)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = :id", namedStmt)

	_, err := r.db.NamedExecContext(ctx, query, userDto)
	if err != nil {
		log.Error("[UserRepository][Update] an error occurred: ", err)
		return err
	}

	return nil
}

// IncrementLoginCount increments a user's login count by 1
func (r *UserRepository) IncrementLoginCount(ctx context.Context, userID string) error {
	query := fmt.Sprintf("UPDATE users SET login_count = login_count + 1 WHERE id = $1")
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Error("[UserRepository][IncrementLoginCount] an error occurred: ", err)
	}

	return nil
}
