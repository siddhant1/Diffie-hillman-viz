package user

import (
	"database/sql"
	"e2e-api-server/db/wrappers"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
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

	user := wrappers.UserRequest{}

	c.Bind(&user)

	insertStmt := `INSERT INTO users (username, password) VALUES ($1, $2)`
	fmt.Print(user)
	password, err := HashPassword(user.Password)
	CheckError(err)

	_, err = db.Exec(insertStmt, user.Username, password)
	CheckError(err)

	return c.JSON(http.StatusCreated, user.Username)

}

func SyncPublicKey(c echo.Context) error {
	db := c.Get("postgres").(*sql.DB)

	user := wrappers.PublicKeyRequest{}
	c.Bind((&user))

	rows, err := db.Query("Select count(*) from KEYS where name = $1", user.Name)

	CheckError(err)
	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	if count == 0 {

		insertStmt := `INSERT INTO KEYS (name, keys) VALUES ($1, $2)`
		_, err = db.Exec(insertStmt, user.Name, user.PublicKey)
		CheckError(err)

	}
	rows, err = db.Query("Select name, keys from KEYS")
	CheckError(err)

	resultMap := make(map[string]string)

	for rows.Next() {
		var id string
		var publicKey string
		err := rows.Scan(&id, &publicKey)
		if err != nil {
			log.Fatal(err)
		}
		resultMap[id] = publicKey
	}

	return c.JSON(http.StatusOK, resultMap)

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
