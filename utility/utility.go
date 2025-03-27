package utility

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"example.com/try-echo/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func StringToHMACSHA256(message, secret string) string {
	hmacHash := hmac.New(sha256.New, []byte(secret))
	hmacHash.Write([]byte(message))
	hash := hmacHash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}

func GenerateLoginJWT(userUniqueId string, userFullName string, userEmail string, userIdRole string) (string, error) {
	conf := config.GetConfig()

	claims := jwt.MapClaims{
		"userUniqueId": userUniqueId,
		"userFullName": userFullName,
		"userEmail":    userEmail,
		"userIdRole":   userIdRole,
		"exp":          time.Now().Add(time.Minute * time.Duration(conf.JWT_TIME_MINUTES)).Unix(), // Token expiration time (72 hours)
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString([]byte(conf.JWT_SECRET))
	if err != nil {
		return err.Error(), err
	}

	return tokenString, nil
}

func ExtractJWTToken(c echo.Context) (string, error) {
	// Get the Authorization header from the request
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header does not start with Bearer")
	}

	// Extract the token part (after "Bearer ")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}
