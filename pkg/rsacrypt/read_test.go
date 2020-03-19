package rsacrypt

import "testing"

func TestGetPublickKeys(t *testing.T) {
	file := "../rsacrypt_test/id_rsa_test.pub"
	if _, err := GetPublickKeys(file); err != nil {
		t.Errorf("get public key failed, err:%v\n", err)
	}
}
