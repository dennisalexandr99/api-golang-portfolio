package controllers

import (
	"net/http"

	"strconv"

	"example.com/try-echo/models"
	"github.com/labstack/echo"
)

func FetchAllUser(c echo.Context) error {
	result, err := models.FetchAllUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func CreateNewUser(c echo.Context) error {
	newFullName := c.FormValue("newFullName")
	newUniqueId := c.FormValue("newUniqueId")
	newEmail := c.FormValue("newEmail")
	newPassword := c.FormValue("newPassword")
	newIdRole := c.FormValue("newIdRole")
	intNewIdRole, err := strconv.Atoi(newIdRole)

	result, err := models.CreateNewUser(newFullName, newUniqueId, newEmail, newPassword, intNewIdRole)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func EditUser(c echo.Context) error {
	targetUserUniqueId := c.FormValue("targetUserUniqueId")
	newFullName := c.FormValue("newFullName")
	newEmail := c.FormValue("newEmail")
	newPassword := c.FormValue("newPassword")
	newIdRole := c.FormValue("newIdRole")
	intNewIdRole, err := strconv.Atoi(newIdRole)

	result, err := models.EditUser(targetUserUniqueId, newFullName, newEmail, newPassword, intNewIdRole)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteUser(c echo.Context) error {
	userUniqueId := c.FormValue("userUniqueId")
	password := c.FormValue("password")

	result, err := models.DeleteUser(userUniqueId, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
