package aescrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/randomStr"
)

// type AesKey []byte
type AesKey struct {
	Key []byte
}

// CreateKey
// create random AES key
func CreateRandomKey() (key *AesKey) {
	newkey := []byte(randomStr.GetRandomStr(32))

	aesKey := AesKey{newkey}
	return &aesKey
}

func GetKey(file string) (key *AesKey) {
	thekey, err := ioutil.ReadFile(file) //Check to see if a key was already created
	if err != nil {
		key = CreateRandomKey() //If not, create one
	} else {
		key = &AesKey{thekey} //If so, set key as the key found in the file
	}
	return
}

// EncryptFile
func (key *AesKey) EncryptFile(inputfile string, outputfile string) {
	b, err := ioutil.ReadFile(inputfile) //Read the target file
	if err != nil {
		logs.Debug("Unable to open the input file!\n")
		os.Exit(0)
	}
	ciphertext := key.encrypt(b)
	//logs.Debug("%x\n", ciphertext)
	err = ioutil.WriteFile(outputfile, ciphertext, 0644)
	if err != nil {
		logs.Debug("Unable to create encrypted file!\n")
		os.Exit(0)
	}
}

func (key *AesKey) DecryptFile(inputfile string, decryptedFile string) (err error) {
	if "" == decryptedFile {
		logs.Debug("decryptedFile can't be empty!")
	}
	z, err := ioutil.ReadFile(inputfile)
	if err != nil {
		err = fmt.Errorf("Unable to Read decrypted file:%v!, err:%v\n", inputfile, err)
		logs.Error(err.Error())
		return
	}
	result := key.decrypt(z)
	//logs.Debug("Decrypted: %s\n", result)
	logs.Debug("Decrypted file:%v was created with file permissions 0777\n", decryptedFile)
	err = ioutil.WriteFile(decryptedFile, result, 0777)
	if err != nil {
		logs.Debug("Unable to create decrypted file!, err:%v\n", err)
		os.Exit(0)
	}
	return
}

func encodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func decodeBase64(b []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		logs.Debug("Error: Bad Key!\n")
		os.Exit(0)
	}
	return data
}

func (key *AesKey) encrypt(text []byte) []byte {
	block, err := aes.NewCipher(key.Key)
	if err != nil {
		logs.Debug("Error, NewCipher, err:%v\n", err)
		panic(err)
	}
	b := encodeBase64(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return ciphertext
}
func (key *AesKey) decrypt(text []byte) []byte {
	block, err := aes.NewCipher(key.Key)
	if err != nil {
		panic(err)
	}
	if len(text) < aes.BlockSize {
		logs.Debug("Error!\n")
		os.Exit(0)
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return decodeBase64(text)
}
