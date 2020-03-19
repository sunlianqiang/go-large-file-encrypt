package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/compress"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/rsacrypt"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/s3"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/sha"
)

// Command-line flags
var (
	mode    = flag.String("mode", "encrypt", "Encrypt file using public key; or decrypt using private key")
	infile  = flag.String("in", "", "Path to input file")
	outfile = flag.String("out", "", "Path to output file")

	// encrypt
	outkeyfile = flag.String("outkey", "", "Option. Path to random AES key")
	pubKey     = flag.String("pubkey", "~/.ssh/id_rsa.pub", "Required. Path to RSA public key for decryption")
	s3Bucket   = flag.String("s3_bucket", "agora-finace-backup-test", "Required. aws s3 bucket")
	s3Region   = flag.String("s3_region", "cn-north-1", "Required. aws s3 region")

	// decrypt
	privateKey = flag.String("privatekey", "~/.ssh/id_rsa", "Required. Path to RSA private key for decryption")
	aesKeyfile = flag.String("aeskey", "", "AES key encrypted by RSA public key. default is [infile].aeskey-0.enc")
)

func clean(files []string) {
	for _, v := range files {
		err := os.Remove(v)

		if err != nil {
			logs.Debug("remove file:%v, err:%v\n", v, err)
			return
		}
		logs.Debug("File:%v successfully deleted\n", v)
	}
}

func encrypt() {
	if *infile == "" || *pubKey == "" {
		flag.PrintDefaults()
		return
	}
	infileName := filepath.Base(*infile)

	// 源文件AES加密后文件
	if "" == *outfile {
		*outfile = infileName + ".enc"
	}
	// AES加密后文件
	if "" == *outkeyfile {
		*outkeyfile = infileName + ".key.enc"
	}
	logs.Debug("Encrypting file:%v, outfile:%v, out aes key file:%v\n", *infile, *outfile, *outkeyfile)

	// read in public key
	// *pubKey = "../../pkg/rsacrypt_test/id_rsa_test.pub"
	keyStrMap, err := rsacrypt.GetPublickKeys(*pubKey)
	if err != nil {
		log.Println("Unable to read public key file. ", err)
	}

	// encrypt
	timeLast := time.Now()
	var outkeyFileArr []string
	if err = rsacrypt.Encrypt(&keyStrMap, *infile, *outfile, &outkeyFileArr); err != nil {
		log.Fatalf(": %s", err)
		return
	}
	logs.Debug("use time:%v, encrypt \n", time.Since(timeLast))
	timeLast = time.Now()
	// get sha256 of original file
	hashFile, err := sha.CreateSha256File(*infile)
	if err != nil {
		log.Fatalf("CreateSha256File error : %s", err)
		return
	}

	// package encrypted file, encrypted AES key file, sha256 file

	files := []string{*outfile, hashFile}
	files = append(files, outkeyFileArr...)
	zipFile := infileName + "." + time.Now().Local().Format(time.RFC3339) + ".zip"

	if err := compress.ZipFiles(zipFile, files); err != nil {
		panic(err)
		return
	}
	logs.Debug("use time:%v, Zipped File:%v, \n", time.Since(timeLast), zipFile)
	timeLast = time.Now()

	// upload to s3
	runtime.GOMAXPROCS(runtime.NumCPU())
	// S3_REGION := "cn-north-1"
	// S3_BUCKET := "agora-finace-backup-test"
	logs.Debug("aws s3 region:%v, bucket:%v, upload zip file:%v\n", *s3Region, *s3Bucket, zipFile)
	if err := s3.UploadMultipart(*s3Region, *s3Bucket, zipFile); err != nil {
		errStr := fmt.Sprintf("Fail to upload file:%v to s3, region:%v, bucket:%v, err:%v", *&zipFile, *s3Region, *s3Bucket, err)
		fmt.Errorf(errStr)
		return
	}
	logs.Debug("use time:%v, upload file:%v to s3, region:%v, bucket:%v", time.Since(timeLast), *&zipFile, *s3Region, *s3Bucket)

	logs.Debug("start clean temp files\n")
	files = append(files, zipFile)
	clean(files)
}

func decrypt() {
	if *privateKey == "" || *infile == "" || *outfile == "" {
		flag.PrintDefaults()
		return
	}

	if "" == *aesKeyfile {
		*aesKeyfile = *infile + ".aeskey-0.enc"
	}
	// read in private key
	// privateKey = "../../pkg/rsacrypt_test/id_rsa_test"
	keybytes, err := ioutil.ReadFile(*privateKey)
	if err != nil {
		log.Fatalf("Unable to read private key file. %s", err)
	}

	if "" == *outfile {
		*outfile = *infile + ".dec"
	}
	// decrypt
	// aesKeyfile := (*infile)[0:len(*infile)-4] + ".key.enc"

	if err := rsacrypt.Decrypt(keybytes, *infile, *aesKeyfile, *outfile); err != nil {
		log.Fatalf("Decrypt error : %s", err)
	}
	log.Println("Decryption ..done")
}
func main() {
	runtime.GOMAXPROCS(1)

	flag.Parse()

	switch *mode {
	case "encrypt":
		encrypt()
	case "decrypt":
		decrypt()
	default:
		log.Fatal("Unknown mode. Valid option is encrypt or decrypt ")
	}

}
