package main

import (
	"database/sql"
	"e2e-api-server/db"
	handlers "e2e-api-server/handlers"
	"e2e-api-server/handlers/user"
	rngprime "e2e-api-server/rngPrime"

	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
)

func primeKeyMiddleWare(pk rngprime.PrimeKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("prime-numbers", pk)
			return next(c)
		}
	}
}

func dbMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("postgres", db)
			return next(c)
		}
	}
}

func main() {

	PrimeKey := rngprime.GeneratePublicNumbers()

	e := echo.New()
	db := db.NewPostgresClient()
	e.Use(primeKeyMiddleWare(PrimeKey))
	e.Use(dbMiddleware(db))

	server := socketio.NewServer(nil)

	server.OnEvent("/", "share-key", func(s socketio.Conn, msg string) {
		s.Emit("reply", "have "+msg)
	})

	e.GET("/api/get_prime_numbers", handlers.GetPublicPrimeNumbers)
	e.POST("/api/login", user.LoginUser)
	e.POST("/api/signup", user.SignUpUser)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
