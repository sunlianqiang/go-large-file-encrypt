package gpg

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestEnc(t *testing.T) {
	// EncryptPGP(msg2Enc string, pubkeyFile string, outfile string)
	msg2Enc := "this is a aes secret"
	pubkeyFile := "slq-pub-A12C6.asc"
	outfile := "aeskey.gpg"
	EncryptPGP(msg2Enc, pubkeyFile, outfile)

}

func TestDec(t *testing.T) {
	testFile := "aeskey-slq-pub-A12C6.asc-0.enc"
	aeskey, err := DecryptGPG(testFile)
	logs.Debug("aeskey:%s, err:%v", aeskey, err)

}
