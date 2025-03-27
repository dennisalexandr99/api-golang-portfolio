package routes

import (
	"net/http"

	"example.com/try-echo/controllers"
	"example.com/try-echo/middleware"
	"github.com/labstack/echo"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello from echo")
	})

	e.GET("/login", controllers.Login)
	e.GET("/user", controllers.FetchAllUser, middleware.IsAuthenticated)
	e.POST("/user", controllers.CreateNewUser, middleware.IsAuthenticated)
	e.PUT("/user", controllers.EditUser, middleware.IsAuthenticated)
	e.DELETE("/user", controllers.DeleteUser, middleware.IsAuthenticated)

	// e.Logger.Fatal(e.Start(":6969"))
	return e
}
