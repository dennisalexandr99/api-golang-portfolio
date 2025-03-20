package routes

import (
	"net/http"

	"example.com/try-echo/controllers"
	"github.com/labstack/echo"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello from echo")
	})

	e.GET("/user", controllers.FetchAllUser)
	e.POST("/user", controllers.CreateNewUser)
	e.PUT("/user", controllers.EditUser)
	e.DELETE("/user", controllers.DeleteUser)

	// e.Logger.Fatal(e.Start(":6969"))
	return e
}
