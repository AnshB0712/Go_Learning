package computer

type Computer struct { // exported type
	Brand string  // exported field
	Price float64 // exported field
	year  int     // unexported field
}