package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

// jwt service
type JWTService interface {
	GenerateToken(UserNumber string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
type AuthCustomClaims struct {
	UserNumber string `json:"user_number"`
	jwt.RegisteredClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

// auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Hamsoft",
	}
}

// getSecretKey returns the configured JWT signing secret.
//
// Release 1 (EMS security fix B1): the previous "secret" fallback was removed.
// A missing JWT_SECRET must never silently produce a predictable signing key.
// Startup is refused in main.go when JWT_SECRET is empty.
func getSecretKey() string {
	return os.Getenv("JWT_SECRET")
}

func (service *jwtServices) GenerateToken(UserNumber string) (string, error) {
	claims := &AuthCustomClaims{
		UserNumber,
		jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48).Unix(11323,2)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    service.issure,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))

	return t, err
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(encodedToken, &AuthCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})

}
