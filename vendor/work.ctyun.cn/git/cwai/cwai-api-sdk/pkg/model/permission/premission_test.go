package permission

import "testing"

func TestSm4Decrypt(t *testing.T) {
	sk := ""
	desk := Sm4Decrypt(sk)
	println(desk)
}
