package rsacrypt_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sunlianqiang/go-large-file-encrypt/pkg/rsacrypt"
)

const (
	PubKeyFile    = "./id_rsa_test.pub"
	PrivKeyFile   = "./id_rsa_test"
	InFile        = "./test.txt"
	OutFile       = InFile + ".enc"
	OutKeyFile    = InFile + ".key.enc"
	DecryptedFile = "./test-decrypted.txt"
)

func TestEncrypt_Success(t *testing.T) {
	keybytes, err := ioutil.ReadFile(PubKeyFile)
	if err != nil {
		t.Error("Unable to read public key file. ", err)
	}
	err = rsacrypt.Encrypt(keybytes, InFile, OutFile, OutKeyFile)
	if err != nil {
		t.Error("TestEncrypt Failed. ", err)
	}
	f1, _ := ioutil.ReadFile(InFile)
	f2, _ := ioutil.ReadFile(OutFile)
	if bytes.Equal(f1, f2) {
		t.Fail()
	}
}

func TestEncrypt_InvalidKey(t *testing.T) {
	err := rsacrypt.Encrypt([]byte("invalid_key"), InFile, OutFile, OutKeyFile)
	if err == nil {
		t.Fail()
	}
}

func TestEncrypt_InvalidInfile(t *testing.T) {
	keybytes, err := ioutil.ReadFile(PubKeyFile)
	err = rsacrypt.Encrypt(keybytes, "invalid_infile", OutFile, OutKeyFile)
	if err == nil {
		t.Fail()
	}
}

func TestDecrypt_Success(t *testing.T) {
	defer cleanup()
	keybytes, err := ioutil.ReadFile(PrivKeyFile)
	if err != nil {
		t.Error("Unable to read private key file. ", err)
	}
	err = rsacrypt.Decrypt(keybytes, OutFile, OutKeyFile, DecryptedFile)
	if err != nil {
		t.Error("TestDecrypt Failed. ", err)
	}
	f1, _ := ioutil.ReadFile(InFile)
	f2, _ := ioutil.ReadFile(DecryptedFile)
	if !bytes.Equal(f1, f2) {
		t.Fail()
	}
}

func cleanup() {
	os.Remove(OutFile)
	os.Remove(OutKeyFile)
	os.Remove(DecryptedFile)
}
