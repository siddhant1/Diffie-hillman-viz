package rngprime

import (
	"crypto/rand"
)

type PrimeKey struct {
	G int32  `json:"id,string"`
	N string `json:"n,string"`
}

func GeneratePublicNumbers() PrimeKey {
	var MAX_BITS int = 2000

	for {
		// Generate a random number with the specified bit size
		num, _ := rand.Prime(rand.Reader, MAX_BITS)

		return PrimeKey{
			G: 17,
			N: num.Text(10),
		}
	}
}
