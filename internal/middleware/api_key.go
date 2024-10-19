package middleware

import (
	"SkinRest/config"
	"SkinRest/internal/database"
	"SkinRest/pkg/models"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeader = "Authorization"
)

func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		appctx, exists := c.MustGet("appCtx").(*database.AppContext)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader(authorizationHeader)

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrTokenNotProvided.Error()})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidTokenFormat.Error()})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userData, err := appctx.GetUserFromToken(token)
		if err != nil {
			if err.Error() == models.ErrUserNotFound.Error() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			appctx.Logger.Error(err.Error())
			c.Abort()
			return

		}

		c.Set("userData", userData)
		c.Next()

	}
}

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader(authorizationHeader)

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrTokenNotProvided.Error()})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidTokenFormat.Error()})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidSigningMethod.Error()})
				c.Abort()
				return nil, nil
			}
			return []byte(config.GetConfig().Auth.JwtSecret), nil
		})

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidToken.Error() + ": " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		claims := token.Claims

		expTime, err := claims.GetExpirationTime()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidTokenClaims.Error() + ": 'exp' claim is missing or not an integer."})
			c.Abort()
			return
		}

		if expTime.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrTokenExpired.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
