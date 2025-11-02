package api

import (
	"net/http"

	"github.com/crayon/wrap-bot/internal/admin/middleware"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Username == "admin" && req.Password == "admin" {
		token, err := middleware.GenerateToken(req.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		}
		return c.JSON(http.StatusOK, LoginResponse{Token: token})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
}
