package log

import (
	"fmt"
	"github.com/woshikedayaa/news-poster/pkg/utils/futil"
	"os"
	"runtime"
	"testing"
)

func TestCaller(t *testing.T) {
	fmt.Println(runtime.Caller(0))
	fmt.Println(os.Getwd())
	fmt.Println(os.Args[0])
	fmt.Println(futil.GetRootDir())
}
