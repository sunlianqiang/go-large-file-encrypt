package aescrypt

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestAes(t *testing.T) {

	key := CreateRandomKey()
	logs.Debug("Key: %x\n", *key)
	// GetAesRandomSecrets()
	// encrypt original file with AES
	var originalFile, encryptedFile, decryptedFile string
	// originalFileName = "/Users/slq/Downloads/Apollo.11.2019.1080p.BluRay.x264.DTS-HD.MA.5.1-FGT/Apollo.11.2019.1080p.BluRay.x264.DTS-HD.MA.5.1-FGT.mkv"
	originalFile = "origin.txt"
	encryptedFile = "encryptedfile"
	decryptedFile = "decryptedFile"
	key.EncryptFile(originalFile, encryptedFile)

	key.DecryptFile(encryptedFile, decryptedFile)
}
