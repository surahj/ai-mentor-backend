package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/surahj/ai-mentor-backend/app/library"
	"github.com/surahj/ai-mentor-backend/app/models"
)

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid token",
			})
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid user_id in token",
			})
		}
		userID := int64(userIDFloat)
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid user_id in token",
			})
		}
		c.Set("user_id", userID)

		log.Printf("User ID: %d", userID)

		_, err = library.GetUserByID(userID)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid token",
			})
		}

		return next(c)
	}
}
