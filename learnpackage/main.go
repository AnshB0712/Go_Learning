package main

import (
	"fmt"
	"learnpackage/simpleinterest"
	"log"
)

var p,r,t = -8000,3.25,1

func init(){
	fmt.Println("package: main initialised")
	if p < 0 {
		log.Fatal("Principal is less than zero")
	}
	if r < 0 {
		log.Fatal("Rate of interest is less than zero")
	}
	if t < 0 {
		log.Fatal("Duration is less than zero")
	}
}

func main() {
	fmt.Println("Simple interest calculation")
	s := simpleinterest.SimpleI(8000, 3.25, 1)
	fmt.Printf("The Simple Interset is %f", s)
}