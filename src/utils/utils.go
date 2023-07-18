// Utils includes all the utility functions required in this project
package utils

// Max returns the maximum value between 2 integers.
//
// It takes two integers of type int64.
//
// Parameters:
// - x: the first integer
// - y: the second integer
//
// Returns:
// - the max of 2 input as int64
//
// Example:
//
//	maxi := Max(10, 20)
//	fmt.Println(maxi)
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
