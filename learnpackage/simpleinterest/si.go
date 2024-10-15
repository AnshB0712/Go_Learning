package simpleinterest

import "fmt"

func init() {
	fmt.Println("package: simpleinterest is initialised")
}

func SimpleI(p, r, t float64) float64 {
	return p * r * t / 100
}