package utils

import (
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/huseyinozsoy/go-jwt/models"
)

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(uuid uuid.UUID) string {
	_, accessToken := GenerateAccessClaims(uuid)
	return accessToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(uuid uuid.UUID) (*models.Claims, string) {

	t := time.Now()
	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid.String(),
			ExpiresAt: t.Add(30 * time.Minute).Unix(),
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
		ID: uuid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			c.SendStatus(fiber.StatusUnauthorized)
			return nil
		}

		// Extract the token from the "Bearer <token>" format
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessToken == authHeader {
			c.SendStatus(fiber.StatusUnauthorized)
			return nil
		}
		claims := new(models.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"general": "Token Expired",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// this is not even a token, we should delete the cookies here
				c.ClearCookie("access_token", "refresh_token")
				c.SendStatus(fiber.StatusForbidden)
				return nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				c.SendStatus(fiber.StatusUnauthorized)
				return nil
			} else {
				// cannot handle this token
				c.ClearCookie("access_token", "refresh_token")
				c.SendStatus(fiber.StatusForbidden)
				return nil
			}
		}

		c.Locals("id", claims.Issuer)
		c.Next()

		return nil
	}
}
