package models

import (
	"net/http"

	"fmt"
	"strconv"

	"example.com/try-echo/config"
	"example.com/try-echo/db"
	"example.com/try-echo/utility"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type User struct {
	IdUser       int    `json:"id_user"`
	UserUniqueId string `json:"user_unique_id"`
	UserFullname string `json:"user_fullname"`
	UserEmail    string `json:"user_email"`
	IdRole       string `json:"id_role"`
}

type UserLogin struct {
	JWTToken     string `json:"jwt_token"`
	JWTExpires   int    `json:"jwt_expires"`
	IdUser       int    `json:"id_user"`
	UserUniqueId string `json:"user_unique_id"`
	UserFullname string `json:"user_fullname"`
	UserEmail    string `json:"user_email"`
	IdRole       string `json:"id_role"`
}

func Login(userUniqueId string, password string) (Response, error) {
	var obj UserLogin
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()
	hashedPassword := utility.StringToHMACSHA256(password, conf.HMAC256_SECRET)

	rows, err := con.Query("select id_user, user_unique_id, user_fullname, user_email, id_role from user where user_unique_id = ? AND user_password = ?", userUniqueId, hashedPassword)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&obj.IdUser, &obj.UserUniqueId, &obj.UserFullname, &obj.UserEmail, &obj.IdRole)
		if err != nil {
			return res, err
		}

		token, err := utility.GenerateLoginJWT(obj.UserUniqueId, obj.UserFullname, obj.UserEmail, obj.IdRole)
		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Message = "Failed to generate token"
		}

		obj.JWTToken = token
		obj.JWTExpires = conf.JWT_TIME_MINUTES

		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = obj
	} else {
		res.Status = http.StatusNotFound
		res.Message = "Wrong username/password"
	}

	return res, nil
}

func FetchAllUser() (Response, error) {
	var obj User
	var arrobj []User
	var res Response
	con := db.CreateCon()

	sqlStatement := "select id_user, user_unique_id, user_fullname, user_email, id_role from user"

	rows, err := con.Query(sqlStatement)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.IdUser, &obj.UserUniqueId, &obj.UserFullname, &obj.UserEmail, &obj.IdRole)
		if err != nil {
			return res, err
		}

		arrobj = append(arrobj, obj)
	}
	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrobj

	return res, nil
}

func CreateNewUser(c echo.Context, newFullName string, newUniqueId string, newEmail string, newPassword string, newIdRole int) (Response, error) {
	var obj User
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()

	tokenString, err := utility.ExtractJWTToken(c)
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWT_SECRET), nil
	})

	userIdRoleInt, err := strconv.Atoi(claims["userIdRole"].(string))

	if userIdRoleInt < 2 {
		res.Status = http.StatusForbidden
		res.Message = "Your role isn't capable to create new account"
		return res, nil
	}

	hashedPassword := utility.StringToHMACSHA256(newPassword, conf.HMAC256_SECRET)
	result, err := con.Exec("INSERT INTO user (user_unique_id, user_fullname, user_email, user_password, id_role) VALUES (?,?,?,?,?)", newUniqueId, newFullName, newEmail, hashedPassword, newIdRole)
	if err != nil {
		return res, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return res, err
	}

	rows, err := con.Query("SELECT id_user, user_unique_id, user_fullname, user_email, id_role FROM user WHERE id_user = ?", lastInsertedId)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&obj.IdUser, &obj.UserUniqueId, &obj.UserFullname, &obj.UserEmail, &obj.IdRole)
		if err != nil {
			return res, err
		}

		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = obj
	} else {
		res.Status = http.StatusBadRequest
		res.Message = "Failed to create a new user"
	}

	return res, nil
}

func EditUser(c echo.Context, newFullName string, newEmail string, newPassword string, newIdRole int) (Response, error) {
	var obj User
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()

	tokenString, err := utility.ExtractJWTToken(c)
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWT_SECRET), nil
	})

	hashedPassword := utility.StringToHMACSHA256(newPassword, conf.HMAC256_SECRET)
	result, err := con.Exec("UPDATE user SET user_fullname = ?, user_email = ?, user_password = ?, id_role = ? WHERE user_unique_id = ?;", newFullName, newEmail, hashedPassword, newIdRole, claims["userUniqueId"].(string))
	if err != nil {
		return res, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return res, fmt.Errorf("no rows updated")
	}

	rows, err := con.Query("SELECT id_user, user_unique_id, user_fullname, user_email, id_role FROM user WHERE user_unique_id = ?", claims["userUniqueId"].(string))
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&obj.IdUser, &obj.UserUniqueId, &obj.UserFullname, &obj.UserEmail, &obj.IdRole)
		if err != nil {
			return res, err
		}

		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = obj
	} else {
		res.Status = http.StatusBadGateway
		res.Message = "Failed to edit user"
	}

	return res, nil
}

func DeleteUser(c echo.Context, password string) (Response, error) {
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()

	tokenString, err := utility.ExtractJWTToken(c)
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWT_SECRET), nil
	})

	hashedPassword := utility.StringToHMACSHA256(password, conf.HMAC256_SECRET)
	result, err := con.Exec("DELETE from user where user_unique_id = ? AND user_password = ?;", claims["userUniqueId"].(string), hashedPassword)
	if err != nil {
		return res, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return res, fmt.Errorf("no rows updated")
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = "Rows affected: " + strconv.FormatInt(rowsAffected, 10)

	return res, nil
}
