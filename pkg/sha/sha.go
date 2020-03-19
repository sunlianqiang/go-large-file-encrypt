package sha

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

func GetSha256(file string) (shaStr string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	shaStr = fmt.Sprintf("%x", h.Sum(nil))
	// logs.Debug("sha256 for file:%v is:%v\n", file, shaStr)
	return
}

func CreateSha256File(inputFile string) (hashFile string, err error) {
	// get sha256 of original file
	shaStr := GetSha256(inputFile)
	logs.Debug("shaStr for infile file:%v, is :%v\n", inputFile, shaStr)
	hashFile = filepath.Base(inputFile) + ".sha256"
	err = ioutil.WriteFile(hashFile, []byte(shaStr), 0644)
	if err != nil {
		logs.Debug("Unable to create encrypted file!\n")
		os.Exit(0)
	}

	return
}
