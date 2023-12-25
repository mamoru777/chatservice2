package jwtservice

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

const (
	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = 30 * 24 * time.Hour
	SecretKey            = "mamoru" // Замените на ваш секретный ключ
)

func CreateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))

	/*token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		Subject:   userId.String(),
	})

	return token.SignedString([]byte(SecretKey))*/
}
func CreateRefreshToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime, expOk := claims["exp"].(float64)
			if !expOk {
				return nil, fmt.Errorf("Invalid token: Missing 'exp' claim")
			}

			// Преобразование времени из float64 в тип time.Time
			expiration := time.Unix(int64(expirationTime), 0)

			// Проверка срока годности токена
			if time.Now().After(expiration) {
				return nil, fmt.Errorf("Token has expired")
			}
			return claims, nil
		}

		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
