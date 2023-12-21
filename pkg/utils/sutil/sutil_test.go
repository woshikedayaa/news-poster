package sutil

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	for i := 0; i < 1000; i++ {
		GenerateRandomString(i)
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomString(1024)
	}
}
