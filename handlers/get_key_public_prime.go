package handlers

import (
	rngprime "e2e-api-server/rngPrime"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetPublicPrimeNumbers(c echo.Context) error {
	primeKeys := c.Get("prime-numbers").(rngprime.PrimeKey)

	return c.JSON(http.StatusOK, primeKeys)
}
