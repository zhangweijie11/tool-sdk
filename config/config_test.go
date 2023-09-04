package config

import (
	"fmt"
	"testing"
)

func TestGenerates(t *testing.T) {
	var key string
	key = "YunzhiSec"
	fmt.Println(key)
	// elastic pass
	str0 := EncryptString("123456", key)
	fmt.Println(str0)
	// mysql pass
	str1 := EncryptString("123456", key)
	fmt.Println(str1)
	// redis pass
	str2 := EncryptString("123456", key)
	fmt.Println(str2)
}
