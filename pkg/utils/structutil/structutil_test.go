package structutil

import (
	"fmt"
	"testing"
)

type cloneTest struct {
	v1 int // private
	V2 int // public
	v3 string
}

func TestClone(t *testing.T) {
	c := &cloneTest{
		v1: 1,
		V2: 2,
		v3: "hello3",
	}
	dst := &cloneTest{}
	Clone(c, dst)
	fmt.Println(dst)
}
