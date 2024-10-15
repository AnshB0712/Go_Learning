package random

import "math/rand"

func GenerateRandomInt() int {
	return rand.Intn(100)
}