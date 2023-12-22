package log

import (
	"fmt"
	"runtime"
	"testing"
)

func TestCaller(t *testing.T) {
	fmt.Println(runtime.Caller(1))
}
