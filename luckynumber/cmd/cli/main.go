package main

import (
	"luckynumber/internal/random"

	"github.com/fatih/color"
)
func main() {
	n := random.GenerateRandomInt()
	c := color.New(color.FgHiMagenta)
	c.Printf("Your Lucky number is %d!!!!\n", n)
}