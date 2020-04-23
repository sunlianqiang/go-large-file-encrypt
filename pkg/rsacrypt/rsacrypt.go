// Package rsacrypt provides methods to encrypt a file and email the
// encrypted file using the users public key and decrypt a file
// using the private key

package rsacrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/aescrypt"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/gpg"
)

// Encrypts the infile using the public key and saves the encrypted
// file as outfile. An error is returned if encryption is unsuccessful.
// infile and outfile should be a pointer type.
func Encrypt(pubKeys []string, infile string, outfile string, outkeyFileArr *[]string) (err error) {
	logs.Info("---------- ENCRYPT START ----------")
	defer logs.Info("---------- ENCRYPT END ----------")
	// create random secrets in file
	aesKey := aescrypt.CreateRandomKey()
	logs.Debug("Success, CreateKey AES Key\n")
	// GetAesRandomSecrets()
	// encrypt original file with AES
	aesKey.EncryptFile(infile, outfile)

	for k, pubKey := range pubKeys {
		// EncryptPGP(msg2Enc, pubkeyFile, outfile)
		outkeyFile := fmt.Sprintf("aeskey-%v-%v.enc", filepath.Base(pubKey), k)
		err = gpg.EncryptPGP(string(aesKey.Key), pubKey, outkeyFile)
		if nil != err {
			fmt.Printf("EncryptPGP err:%v\n", err)
			continue
		}

		*outkeyFileArr = append(*outkeyFileArr, outkeyFile)
	}

	if 0 == len(*outkeyFileArr) {
		errStr := fmt.Sprintf("Fail, no valid public key\n")
		logs.Error(errStr)
		return
	}

	return

}

// Decrypts the infile using the private key keyfile and saves the decrypted
// file as outfile. An error is returned if decryption is unsuccessful.
// keyfile, infile and outfile should be a pointer type.
func Decrypt(keybytes []byte, infile string, inkeyfile string, outfile string) error {

	// decode PEM encoding to ANS.1 PKCS1 DER
	block, _ := pem.Decode(keybytes)
	if block == nil {
		return errors.New("Private key file not PEM-encoded")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return errors.New("Unsupported key type")
	}

	// decode the RSA private key
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	// read the encrypted input key file
	inkeybytes, err := ioutil.ReadFile(inkeyfile)
	if err != nil {
		return fmt.Errorf("Unable to read input key file. %s", err)
	}

	// decrypt the input AES key
	label := []byte("")
	inkey, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, privkey, inkeybytes, label)
	if err != nil {
		return fmt.Errorf("Unable to decrypt AES key. %s", err)
	}

	aesKey := &aescrypt.AesKey{inkey}
	aesKey.DecryptFile(infile, outfile)

	return err
}

// Ecrypt plaintext using AES key.
func encrypt(plaintext []byte, key []byte) ([]byte, error) {

	// create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("Unable to create AES cipher block")
	}

	// empty array of aes block size 16 + plaintext length
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.New("AES initialization vector error")
	}

	// return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, err
}

// Decrypt ciphertext using AES key.
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {

	// create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("Unable to create AES cipher block")
	}

	// check size of ciphertext
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short to decrypt")
	}

	// get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, err
}
