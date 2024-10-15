package main

import "fmt"

func findAvg(n1 int, n2 int) float64 {
	sum := n1 + n2
	avg := float64(sum) / 2
	return avg
}

func calcRec(l float64, b float64)(float64, float64) {
	a := l*b
	p := 2*(l+b)

	return a,p
}

func main() {
	a := findAvg(5, 7)
	fmt.Printf("The Average is %f\n", a)
	a,p := calcRec(5,8)
	fmt.Printf("The Perimeter is %f and The Area is %f\n", a, p)
}