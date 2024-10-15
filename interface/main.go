package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
)

type Book struct {
	name string
	author string
}
func (b Book) String() string {
	return fmt.Sprintf("Reading %v by %v", b.name, b.author)
}

type Coin int
func (c Coin) String() string {
	return fmt.Sprintf("Coin of value %d", c)
}

func WriteLog(l fmt.Stringer) {
	fmt.Println(l.String())
}

type Customer struct {
	Name string
	Age int
}

func (c *Customer) WriteJSON(w io.Writer) error {
	json, err := json.Marshal(c)

	if(err != nil) {
		return err
	}

	_, err = w.Write(json)

	return err
}	

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width float64
	Height float64
}
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (r Rectangle) Perimeter() float64 {
	return 2*r.Width + 2*r.Height
}

type Square struct {
	side float64
}
func (r Square) Area() float64 {
	return r.side * r.side
}
func (r Square) Perimeter() float64 {
	return 2*r.side + 2*r.side
}

type RightTriangle struct {
	base float64
	height float64
}
func (r RightTriangle) Area() float64 {
	return 0.5*r.base*r.height
}
func (r RightTriangle) Perimeter() float64 {
	h := math.Sqrt(r.base*r.base + r.height*r.height)
	return r.base + r.height + h
}

func RandomShape() Shape {
	n := rand.Intn(9)
	if n >= 0 && n < 3 {
		return Rectangle{Width: 10, Height: 5}
	} else if n >= 3 && n < 6 {
		return Square{side: 10}
	} else {
		return RightTriangle{base: 10, height: 5}
	}
}	

func findTypeAndLogProps(s Shape) {
	switch s.(type) {
	case Rectangle:
		fmt.Printf("The area and perimeter of rectangle with sides %f and %f are %f and %f respectively.", s.(Rectangle).Height,s.(Rectangle).Width,s.Area(), s.Perimeter())
	case Square:
		fmt.Printf("The area and perimeter of square with side %f are %f and %f respectively.", s.(Square).side, s.Area(), s.Perimeter())
	case RightTriangle:
		fmt.Printf("The area and perimeter of right triangle with base %f and height %f are %f and %f respectively.", s.(RightTriangle).base,s.(RightTriangle).height,s.Area(), s.Perimeter())	
	default: 
		fmt.Println("Unknown shape")
	}

}
	

func main() {
	b := Book{
		name: "Dune",
		author: "Frank Herbert",
	}
	c := Coin(10)

	WriteLog(b)

	WriteLog(c)

	cust := Customer{
		Name: "Alice",
		Age: 25,
	}
	addressOfCustVar := &cust

	var buf bytes.Buffer
	err := addressOfCustVar.WriteJSON(&buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())

	f, err := os.Create("./customer.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	e := addressOfCustVar.WriteJSON(f)
	if e != nil {
		log.Fatal(e)
	}

	// as of now only reason to use type assertion is if we are getting a value from some exteral source and return type of that call is interface of some kind so we need to type assert it to get the actual value to call right properties and methods on it.
	s := RandomShape()
	findTypeAndLogProps(s)

	

}