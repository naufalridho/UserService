package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/SawitProRecruitment/UserService/internal/liberr"
)

const (
	ErrCodeInvalidJwt = "USER_SERVICE-JWT-INVALID"
)

var (
	ErrInvalidJwt = liberr.NewError("Invalid token", ErrCodeInvalidJwt, "")
)

// JWTClaims extends jwt.RegisteredClaims to have extra information such as ClientID
type JWTClaims struct {
	jwt.RegisteredClaims

	ClientID string
}

const (
	jwtExpiryTime = 1 * time.Hour
)

var (
	signingMethod = jwt.SigningMethodHS256
)

// AccessToken holds an access token string and its lifetime in seconds
type AccessToken struct {
	Value    string
	Lifetime time.Duration
}

// GenerateAccessToken generates a new access token
func GenerateAccessToken(clientID, secretKey string) (*AccessToken, error) {
	token := jwt.NewWithClaims(signingMethod, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiryTime)),
		},
		ClientID: clientID,
	})

	val, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Value:    val,
		Lifetime: jwtExpiryTime,
	}, nil
}

// ParseAccessToken validates an access token and parses it into JWTClaims
func ParseAccessToken(accessTokenStr string, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(accessTokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, ErrInvalidJwt
	}

	claims := token.Claims.(*JWTClaims)
	return claims, nil
}

// AccessTokenFromHeader extracts an access token from HTTP header
func AccessTokenFromHeader(header http.Header) string {
	authHeader := header.Get("Authorization")
	return strings.Split(authHeader, " ")[1]
}
