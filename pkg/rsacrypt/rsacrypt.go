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

	ssh "github.com/ianmcmahon/encoding_ssh"

	"github.com/sunlianqiang/go-large-file-encrypt/pkg/aescrypt"
)

// Encrypts the infile using the public key and saves the encrypted
// file as outfile. An error is returned if encryption is unsuccessful.
// infile and outfile should be a pointer type.
func Encrypt(keybytes []byte, infile string, outfile string, outkeyfile string) error {
	// decode string ssh-rsa format to native type
	pubkey, err := ssh.DecodePublicKey(string(keybytes))
	if err != nil {
		return fmt.Errorf("Unable to decode public key. %s", err)
	}

	// create random secrets in file
	aesKey := aescrypt.CreateKey()
	fmt.Printf("Success, CreateKey AES Key\n")
	// GetAesRandomSecrets()
	// encrypt original file with AES
	aesKey.EncryptFile(infile, outfile)

	// encrypt key with rsa publick key
	cipherkey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pubkey.(*rsa.PublicKey), aesKey.Key, []byte(""))
	if err != nil {
		return fmt.Errorf("Unable to encrypt key. %s", err)
	}

	// write key to output file
	//outkeyfile := "c:/temp/" + filepath.Base(infile) + ".key.enc"
	err = ioutil.WriteFile(outkeyfile, cipherkey, 0600)
	if err != nil {
		return fmt.Errorf("Unable to write to output file. %s", err)
	}

	return err
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
	// // read the encrypted input file
	// infilebytes, err := ioutil.ReadFile(infile)
	// if err != nil {
	// 	return fmt.Errorf("Unable to read input file. %s", err)
	// }

	// // decrypt
	// decipher, err := decrypt(infilebytes, inkey)
	// if err != nil {
	// 	return fmt.Errorf("Unable to decrypt. %s", err)
	// }

	// // write to output file
	// err = ioutil.WriteFile(outfile, decipher, 0600)
	// if err != nil {
	// 	return fmt.Errorf("Unable to write to output file. %s", err)
	// }

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
