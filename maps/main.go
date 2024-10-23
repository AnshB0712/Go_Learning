package main

import "fmt"

func main() {
	code := make(map[string]string)

	code["USD"] = "US Dollar"
	code["GBP"] = "Pound Sterling"
	code["EUR"] = "Euro"
	code["INR"] = "Indian Rupee"

	// var nilMap map[string]string

	// nilMap["USD"] = "US Dollar"

	preMap := map[string]string{
		"USD": "US Dollar",
		"GBP": "Pound Sterling",
		"EUR": "Euro",
	}

	fmt.Print(preMap["INR"])

	// to check if key exists in a map or not
	val,ok := preMap["INR"]
	if ok {
		fmt.Print(val)
	}else{
		fmt.Print("not found key\n")
	}

	// to iterate over a map
	for k,v := range preMap {
		fmt.Printf("%v: %v\n",k,v)
	}

	// to delete a key from map
	delete(preMap, "USD")

	fmt.Print(preMap)
}