package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/repositories"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	Repo repositories.UserRepository
}

var secretKey = []byte("your_secret_key") // Change this to a secure secret

// GenerateJWT generates a new JWT token
func (a AuthService) GenerateJWT(userID int64) (string, error) {
	expireAddSecond := time.Second * time.Duration(a.Repo.Base.Config.AuthTokenExpireTime)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expireAddSecond).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseJWT parses the JWT token and returns the claims
func (a AuthService) ParseJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, err
	}
	return token, *claims, nil
}

func (a AuthService) LoginUser(ctx context.Context, email string, password string) (sqlc.User, error) {
	user, err := a.Repo.FindUserByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, err
	}
	if a.Repo.Base.CheckPasswordHash(password, user.Password) {
		return user, nil
	}

	return sqlc.User{}, nil
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return AuthService{Repo: repo}
}
