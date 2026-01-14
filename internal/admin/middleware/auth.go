package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret []byte
var jwtSecretInitialized bool

func getJWTSecret() []byte {
	if !jwtSecretInitialized {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			logger.Error("JWT_SECRET environment variable is required")
			os.Exit(1)
		}
		jwtSecret = []byte(secret)
		jwtSecretInitialized = true
	}
	return jwtSecret
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return getJWTSecret(), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			c.Set("username", claims.Username)
			return next(c)
		}
	}
}
