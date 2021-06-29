package base

import (
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	ip := GetIPByHost("https://api.ashibro.com/transaction", "127.0.0.1")
	fmt.Println(ip)
}
