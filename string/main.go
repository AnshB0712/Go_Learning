package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	str := "hello_hii_bye"

	fmt.Printf("%s\n", str)

	for i := range 4 {
		fmt.Printf("%v\n", str[i])
		fmt.Printf("%x\n", str[i])
	}

	// 	type stringStruct struct {
	//    str unsafe.Pointer
	//    len int
	// 	}	so they dont have cap field so we cannot append anything to a string inshort there is no way to mutate a string in go

	// a := str[0:5]
	// b := len(str)
	// c := make([]byte, 3) we can do all this with string
	// d := append(str, "A") we cannot append to a string

	s := "hello"
	s += "world what "
	s = s + "is going on"

	fmt.Printf("%s\n", s)

	var b strings.Builder

	b.WriteString("hello")
	b.WriteString("world")
	b.WriteString("what is going on")

	fmt.Printf("%s\n", b.String())

	// zero value of a string is empty string not a nil byte slice even though they almost share same structure they cannot be nil;

	var (
		s1 string
		s2 []byte
	)


	fmt.Printf("%v\n", s1 == "")
	fmt.Printf("%v\n", s2 == nil)


	st := "hello world"
	st2 := "hello André"

	for k,v := range st {
		fmt.Printf("%v, %v, %T, %U, %v\n",k,v,v,v, string(v))
	}

	// order to think about strings in go "string" ---> "code points unicode(rune in go terminology)" ---> "only readable byte slice"
	fmt.Printf("%v\n", len(st)) // both st and st2 are of same length but st2 has a special character(é) so it will be counted as 2 bytes bcz len() returns number of bytes in a string pretty evident when you think strings as byte slices
	fmt.Printf("%v\n", len(st2))

	// The correct way to get the number of characters in a string is to convert it to a rune slice and then get the length of the slice.
	fmt.Printf("%v\n", len([]rune(st2)))
	// or use package utf8
	fmt.Printf("%v\n", utf8.RuneCountInString(st2))


}