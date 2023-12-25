package futil

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestGetRootDir(t *testing.T) {
	fmt.Println(runtime.Caller(0))
	fmt.Println(os.Getwd())
	fmt.Println(os.Args[0])
	fmt.Println(GetRootDir())
}
