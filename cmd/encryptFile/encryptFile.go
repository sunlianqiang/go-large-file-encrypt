package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/aescrypt"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/compress"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/gpg"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/rsacrypt"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/s3"
	"github.com/sunlianqiang/go-large-file-encrypt/pkg/sha"
)

type Pubkeys []string

// Command-line flags
var (
	mode    = flag.String("mode", "encrypt", "Encrypt file using public key; or decrypt using private key")
	infile  = flag.String("in", "", "Path to input file")
	outfile = flag.String("out", "", "Path to output file")
	outDir  = flag.String("outDir", ".", "unzip and decrypted dump fileto this dir")

	// encrypt
	outkeyfile = flag.String("outkey", "", "Option. Path to random AES key")
	// pubKey     = flag.String("pubkey", "~/.ssh/id_rsa.pub", "Required. Path to GPG armor public key file")
	s3Bucket = flag.String("s3_bucket", "agora-finace-backup-test", "Required. aws s3 bucket")
	s3Region = flag.String("s3_region", "cn-north-1", "Required. aws s3 region")
	pubKeys  Pubkeys
	zipFile  = flag.String("zipfile", "", "final zip file name upload to aws s3")

	// decrypt
	privateKey    = flag.String("privatekey", "~/.ssh/id_rsa", "Required. Path to RSA private key for decryption")
	aesKeyfile    = flag.String("aeskey", "", "AES key encrypted by RSA public key. default is [infile].aeskey-0.enc")
	sensitiveFile = flag.String("sensitive", "", "sensitive data need be decrypted")

	// test
	noClean = flag.Bool("noclean", false, "don't clean temp file after encrypt")
)

func (i *Pubkeys) String() string {
	return "Required. Path to GPG armor public key file"
}

func (i *Pubkeys) Set(value string) error {
	*i = append(*i, value)
	return nil
}

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
	if *infile == "" {
		logs.Error("NO infile !\n")
		flag.PrintDefaults()
		return
	} else if 0 == len(pubKeys) {
		logs.Error("NO pubkeys\n")
		flag.PrintDefaults()
		return
	}
	infileName := filepath.Base(*infile)

	// 源文件AES加密后文件
	if "" == *outfile {
		*outfile = "sensitivedata-" + infileName + ".enc"
	}

	// // AES加密后文件
	// if "" == *outkeyfile {
	// 	*outkeyfile = infileName + ".key.enc"
	// }
	// logs.Debug("Encrypting file:%v, outfile:%v, out aes key file:%v\n", *infile, *outfile, *outkeyfile)

	// encrypt
	timeLast := time.Now()
	var outkeyFileArr []string
	if err := rsacrypt.Encrypt([]string(pubKeys), *infile, *outfile, &outkeyFileArr); err != nil {
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
	if "" == *zipFile {
		timeLayout := "2006-01-02T15-04-05Z07"
		*zipFile = infileName + "-" + time.Now().Local().Format(timeLayout) + ".zip"
		if "" != *outDir {
			*zipFile = filepath.Join(*outDir, *zipFile)
		}
	}

	// zipFile = "test.zip"
	if err := compress.ZipFiles(*zipFile, files); err != nil {
		panic(err)
		return
	}
	logs.Debug("use time:%v, Zipped File:%v, \n", time.Since(timeLast), *zipFile)
	timeLast = time.Now()

	// upload to s3
	runtime.GOMAXPROCS(runtime.NumCPU())
	// S3_REGION := "cn-north-1"
	// S3_BUCKET := "agora-finace-backup-test"
	logs.Debug("aws s3 region:%v, bucket:%v, upload zip file:%v\n", *s3Region, *s3Bucket, *zipFile)
	retry := 3
	for i := 1; i <= retry; i++ {
		err := s3.UploadMultipart(*s3Region, *s3Bucket, *zipFile)
		if err != nil {
			errStr := fmt.Sprintf("%vth time upload, Fail to upload file:%v to s3, region:%v, bucket:%v", i, *zipFile, *s3Region, *s3Bucket)
			logs.Error(errStr)
		} else {
			logs.Info("%vth time upload, Success to upload file:%v to s3, region:%v, bucket:%v", i, *zipFile, *s3Region, *s3Bucket)
			break
		}

	}

	logs.Debug("use time:%v, upload file:%v to s3, region:%v, bucket:%v", time.Since(timeLast), *zipFile, *s3Region, *s3Bucket)

	if !*noClean {
		logs.Debug("start clean temp files, noClean:%v\n", *noClean)
		// files = append(files, *zipFile)
		clean(files)
	}

}

func decrypt() {

	if *infile == "" {
		logs.Error("NO infile !\n")
		flag.PrintDefaults()
		return
	}
	if "" == *outfile {
		*outfile = *infile + ".dec"
	}

	// unzip
	// outDir := "test"
	files, err := compress.Unzip(*infile, *outDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))


	unzipFiles, err := ioutil.ReadDir(*outDir)
	if err != nil {
		log.Fatal(err)
	}
	var aesKeyBytes []byte
	for _, f := range unzipFiles {
		aesKeyfile := f.Name()
		if strings.HasPrefix(aesKeyfile, "aeskey") {
			logs.Info("decrypt aeskey:%v", aesKeyfile)
			aesKeyfile = filepath.Join(*outDir, aesKeyfile)
			aesKeyBytes, err = gpg.DecryptGPG(aesKeyfile)
			if err != nil {
				errStr := fmt.Sprintf("Unable to read input key file:%s. %s", aesKeyfile, err)
				logs.Error(errStr)
				continue
			}
			// decrypt success
			// logs.Debug("aeskey:%s, stderr output:\n%v", aesKeyBytes, err)
			break
		}
	}

	// gpg -d aeskey-slq-pub-A12C6.asc-0.enc

	// inkeybytes, err := ioutil.ReadFile(*aesKeyfile)
	// if err != nil {
	// 	errStr := fmt.Sprintf("Unable to read input key file. %s", err)
	// 	logs.Error(errStr)
	// 	return
	// }
	// decrypted by aes
	if "" == *sensitiveFile {
		for _, f := range unzipFiles {
			tmpFile := filepath.Join(*outDir, f.Name())
			if strings.Contains(tmpFile, "sensitivedata-") {
				*sensitiveFile = tmpFile
				logs.Info("get sensitiveFile:%v\n", *sensitiveFile)
				break
			}
		}
	}

	logs.Info("DecryptFile sensitiveFile:%v\n", *sensitiveFile)
	aesKey := &aescrypt.AesKey{Key: aesKeyBytes}
	err = aesKey.DecryptFile(*sensitiveFile, *outfile)
	if err != nil {
		errStr := fmt.Sprintf("Failed aesKey.DecryptFile, err:%v", err)
		log.Fatal(errStr)
	}
	log.Println("Decryption ..done")
}
func main() {
	runtime.GOMAXPROCS(1)

	flag.Var(&pubKeys, "pubkey", "Required. Path to GPG armor public key file")
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
