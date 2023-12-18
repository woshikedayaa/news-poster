package etcd

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	fmt.Println(parseEndpoints())
}

func TestRegistryNewService(t *testing.T) {
	register, err := registryNewService("hello", "world")
	if err != nil {
		t.Fatal(err)
	}
	select {}
	register.close()
}
