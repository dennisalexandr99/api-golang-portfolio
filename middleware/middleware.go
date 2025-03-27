package middleware

import (
	"example.com/try-echo/config"
	"github.com/labstack/echo/middleware"
)

var conf = config.GetConfig()

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(conf.JWT_SECRET),
})
