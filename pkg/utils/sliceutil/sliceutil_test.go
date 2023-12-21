package sliceutil

import (
	"fmt"
	"slices"
	"testing"
)

func TestDeleteElement(t *testing.T) {
	var a = make([]int, 64)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}
	// test
	fmt.Println(DeleteElement(a, 10))
	fmt.Println(slices.Delete(a, 10, 11))
}
