package base

import (
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	ip := GetIPByHost("www.baidu.com", "127.0.0.1")
	fmt.Println(ip)
}
