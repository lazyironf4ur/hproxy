package common

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	stack := []int{
		1, 2, 3, 4, 5, 6, 7, 9,
	}
	arr(stack[:])
	fmt.Println(stack)
}

func arr(a []int) {
	(a)[3] = 10
}
