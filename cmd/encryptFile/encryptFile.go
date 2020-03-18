package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sunlianqiang/go-large-file-encrypt/pkg/compress"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/rsacrypt"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/s3"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/sha"
)

// Command-line flags
var (
	mode       = flag.String("mode", "encrypt", "Encrypt file using public key; or decrypt using private key")
	privateKey = flag.String("privatekey", "~/.ssh/id_rsa", "Path to RSA private key for decryption")
	pubKey     = flag.String("pubkey", "~/.ssh/id_rsa.pub", "Path to RSA public key for decryption")
	infile     = flag.String("in", "", "Path to input file")
	outfile    = flag.String("out", "", "Path to output file")
	outkeyfile = flag.String("outkey", "", "Path to random AES key")
	s3Bucket   = flag.String("s3_bucket", "agora-finace-backup-test", "aws s3 bucket")
	s3Region   = flag.String("s3_region", "cn-north-1", "aws s3 region")
	// add option to sign the file
)

func encrypt() {
	if *infile == "" || *pubKey == "" {
		flag.PrintDefaults()
		return
	}
	if "" == *outfile {
		*outfile = filepath.Base(*infile) + ".enc"
	}
	if "" == *outkeyfile {
		*outkeyfile = *infile + ".key.enc"
	}
	fmt.Printf("Encrypting file:%v, outfile:%v, out aes key file:%v\n", *infile, *outfile, *outkeyfile)

	// read in public key
	// *pubKey = "../../pkg/rsacrypt_test/id_rsa_test.pub"
	keybytes, err := ioutil.ReadFile(*pubKey)
	if err != nil {
		log.Println("Unable to read public key file. ", err)
	}

	// encrypt
	timeLast := time.Now()
	if err = rsacrypt.Encrypt(keybytes, *infile, *outfile, *outkeyfile); err != nil {
		log.Fatalf("Encrypt error : %s", err)
		return
	}
	fmt.Printf("use time:%v, encrypt \n", time.Since(timeLast))
	timeLast = time.Now()
	// get sha256 of original file
	hashFile, err := sha.CreateSha256File(*infile)
	if err != nil {
		log.Fatalf("CreateSha256File error : %s", err)
		return
	}

	// package encrypted file, encrypted AES key file, sha256 file
	files := []string{*outfile, *outkeyfile, hashFile}
	zipFile := *outfile + "." + time.Now().Local().Format(time.RFC3339) + ".zip"

	if err := compress.ZipFiles(zipFile, files); err != nil {
		panic(err)
		return
	}
	fmt.Printf("use time:%v, Zipped File:%v, \n", time.Since(timeLast), zipFile)
	timeLast = time.Now()

	// upload to s3
	runtime.GOMAXPROCS(runtime.NumCPU())
	// S3_REGION := "cn-north-1"
	// S3_BUCKET := "agora-finace-backup-test"
	fmt.Printf("aws s3 region:%v, bucket:%v, upload zip file:%v\n", *s3Region, *s3Bucket, zipFile)
	s3.UploadMultipart(*s3Region, *s3Bucket, zipFile)
	fmt.Printf("use time:%v, upload \n", time.Since(timeLast))
}

func decrypt() {
	if *privateKey == "" || *infile == "" || *outfile == "" {
		flag.PrintDefaults()
		return
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
	aesKeyfile := (*infile)[0:len(*infile)-4] + ".key.enc"
	if err := rsacrypt.Decrypt(keybytes, *infile, aesKeyfile, *outfile); err != nil {
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
