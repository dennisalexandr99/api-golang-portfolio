package models

import (
	"database/sql"
	"net/http"

	"fmt"
	"strconv"

	"example.com/try-echo/config"
	"example.com/try-echo/db"
	"example.com/try-echo/utility"
)

type User struct {
	IdUser       int            `json:"id_user"`
	UserUniqueId sql.NullString `json:"user_unique_id"`
	UserFullname sql.NullString `json:"user_fullname"`
	UserEmail    sql.NullString `json:"user_email"`
	IdRole       sql.NullString `json:"id_role"`
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

func CreateNewUser(newFullName string, newUniqueId string, newEmail string, newPassword string, newIdRole int) (Response, error) {
	var obj User
	var arrobj []User
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()

	sqlStatement := "INSERT INTO user (user_unique_id, user_fullname, user_email, user_password, id_role) VALUES (?,?,?,?,?)"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	hashedPassword := utility.StringToHMACSHA256(newPassword, conf.HMAC256_SECRET)

	result, err := stmt.Exec(newUniqueId, newFullName, newEmail, hashedPassword, newIdRole)
	if err != nil {
		return res, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return res, err
	}

	rows, err := con.Query("SELECT id_user, user_unique_id, user_fullname, user_email, id_role FROM user WHERE id_user = ?", lastInsertedId)
	defer rows.Close()
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

func EditUser(targetUserUniqueId string, newFullName string, newEmail string, newPassword string, newIdRole int) (Response, error) {
	var obj User
	var arrobj []User
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()

	sqlStatement := "UPDATE user SET user_fullname = ?, user_email = ?, user_password = ?, id_role = ? WHERE user_unique_id = ?;"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	hashedPassword := utility.StringToHMACSHA256(newPassword, conf.HMAC256_SECRET)

	result, err := stmt.Exec(newFullName, newEmail, hashedPassword, newIdRole, targetUserUniqueId)
	if err != nil {
		return res, err
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return res, fmt.Errorf("no rows updated")
	}

	rows, err := con.Query("SELECT id_user, user_unique_id, user_fullname, user_email, id_role FROM user WHERE user_unique_id = ?", targetUserUniqueId)
	defer rows.Close()
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

func DeleteUser(userUniqueId string, password string) (Response, error) {
	var res Response
	conf := config.GetConfig()
	con := db.CreateCon()
	hashedPassword := utility.StringToHMACSHA256(password, conf.HMAC256_SECRET)

	sqlStatement := "DELETE from user where user_unique_id = ? AND user_password = ?;"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userUniqueId, hashedPassword)
	if err != nil {
		return res, err
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return res, fmt.Errorf("no rows updated")
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = "Rows affected: " + strconv.FormatInt(rowsAffected, 10)

	return res, nil
}
