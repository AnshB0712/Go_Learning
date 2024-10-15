package main

import "fmt"

type player struct {
	name string
	score int
}

func increment(a int) {
	a += 10
}

func incrementViaPointer(a *int) {
	*a += 10

}

func incrementPlayerScore(player *player) {
	player.score += 10

}

func main() {
	a := 10
	b := 10

	fmt.Printf("a: %v\n", a)
	fmt.Printf("b: %v\n", b)

	increment(a)
	incrementViaPointer(&b)


	fmt.Printf("a: %v\n", a)
	fmt.Printf("b: %v\n", b)

	p1 := player{name: "a", score: 10}

	incrementPlayerScore(&p1)

	fmt.Printf("player: %v\n", p1)
}