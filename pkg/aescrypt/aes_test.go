package aescrypt

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {

	key := CreatePrivKey()
	fmt.Printf("Key: %x\n", *key)
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
