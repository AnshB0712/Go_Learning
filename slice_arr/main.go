package main

import "fmt"


func evenNum(nums []int) []int {
	s := []int{}

	for _,v := range nums {
		if v%2==0 {
			s = append(s, v)
		}
	}

	return s
}

func main() {
	var s1 []string
	s2 := []string{}
	s3 := make([]string, 0)

	fmt.Printf("%v,%v,%v \n", len(s1), cap(s1), s1 != nil)
	fmt.Printf("%v,%v,%v \n",len(s2), cap(s2), s2 != nil)
	fmt.Printf("%v,%v,%v \n",len(s3), cap(s3), s3 != nil)

	e := evenNum([]int{1})

	if len(e)==0 {
		fmt.Print("no even founds")
	}

	fmt.Printf("%v\n",e)

	// -------------------- x -------------- //

	arr := [6]int{1,2,3,4,5,6}

	fmt.Printf("original slice: %v\n",arr)

	sliced := arr[1:4]


	fmt.Printf("slice of original %v,%v,%v \n",sliced, len(sliced), cap((sliced)))
	sliced[0] = 66
	fmt.Printf("slice of original mutated %v,%v,%v \n",sliced, len(sliced), cap((sliced)))
	
	fmt.Printf("original slice: %v\n",arr)
	// -------------------- x -------------- //

	// above phenomena "x" happens because of how slice is defined in go -> slice 
	//	struct {
	//     array unsafe.Pointer
	//     len   int
	//     cap   int
	// }

	// to escape "x"

	arr2 := []int{-1,-2,-3,-4,-5,-6}
	sliced2 := make([]int, len((arr2[1:4])))

	fmt.Printf("original arr2 %v\n",arr2)


	fmt.Printf("slice of original %v,%v,%v \n",sliced2, len(sliced2), cap((sliced2)))
	sliced2[0] = 66
	fmt.Printf("slice of original mutated %v,%v,%v \n",sliced2, len(sliced2), cap((sliced2)))
	
	fmt.Printf("original slice: %v\n",arr2)

	arr3 :=[]int{}
	fmt.Printf("original arr3: %v, %v, %v\n", arr3, len(arr3), cap(arr3))
	arr3 = append(arr3, 1)
	fmt.Printf("appendend arr3: %v, %v, %v\n", arr3, len(arr3), cap(arr3))
	arr3 = append(arr3, 1)
	fmt.Printf("appendend arr3: %v,%v, %v\n", arr3, len(arr3), cap(arr3))
	arr3 = append(arr3, 1)
	fmt.Printf("appendend arr3: %v,%v, %v\n", arr3, len(arr3), cap(arr3))
	arr3 = append(arr3, 1)
	fmt.Printf("appendend arr3: %v,%v, %v\n", arr3, len(arr3), cap(arr3))



}