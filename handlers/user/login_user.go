package user

import (
	"database/sql"
	"e2e-api-server/db/wrappers"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const PASSWORD_HASH = "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginUser(c echo.Context) error {
	db := c.Get("postgres").(*sql.DB)
	var user wrappers.UserRequest
	var dbUser wrappers.User
	c.Bind(&user)

	rows, err := db.Query("Select username, password from USERS where username = $1", user.Username)

	CheckError(err)

	for rows.Next() {
		err = rows.Scan(&dbUser.Username, &dbUser.HashedPassword)
		CheckError(err)
	}

	match := CheckPasswordHash(user.Password, PASSWORD_HASH)

	if !match {
		return c.JSON(http.StatusBadRequest, "Wrong Password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": dbUser,
	})

	return c.JSON(http.StatusOK, token)

}

func SignUpUser(c echo.Context) error {
	db := c.Get("postgres").(*sql.DB)
	var user wrappers.UserRequest

	c.Bind(&user)

	insertStmt := `INSERT INTO users (username, password) VALUES ($1, $2)`
	password, err := HashPassword(user.Password)
	CheckError(err)

	_, err = db.Exec(insertStmt, user.Username, password)
	CheckError(err)

	return c.JSON(http.StatusCreated, user)

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
