package database

import (
	"SkinRest/config"
	"SkinRest/pkg/models"
	"log"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateNewToken(user *models.User) string {
	secret := config.GetConfig().Auth.JwtSecret

	payload := &jwt.MapClaims{
		"sub": user.Login,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(14 * 24 * time.Hour).Unix(), // 14 days expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString

}

func CheckToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // check token signing method
			return nil, models.ErrInvalidSigningMethod
		}
		return []byte(config.GetConfig().Auth.JwtSecret), nil
	})

	if err != nil {
		return err
	}

	claims := token.Claims
	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}

	if expTime.Before(time.Now()) { // error if token has expired
		return models.ErrTokenExpired
	}

	return nil
}

func GetPasswordHash(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func ValidatePasswordHash(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
