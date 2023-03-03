package tools

import (
	"fmt"
	"testing"
)

func Test_HashPwd(t *testing.T) {
	pwd, err := PwdHash("qwe123$$")
	if err != nil {
		t.Errorf("hash pwd fail: %v", err)
		return
	}
	fmt.Printf("hash pwd: %s\n", pwd)
}
