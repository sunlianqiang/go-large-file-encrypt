package sha

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	// fmt.Printf("sha256 for file:%v is:%v\n", file, shaStr)
	return
}

func CreateSha256File(inputFile string) (hashFile string, err error) {
	// get sha256 of original file
	shaStr := GetSha256(inputFile)
	fmt.Printf("shaStr for infile file:%v, is :%v\n", inputFile, shaStr)
	hashFile = inputFile + ".sha256"
	err = ioutil.WriteFile(hashFile, []byte(shaStr), 0644)
	if err != nil {
		fmt.Printf("Unable to create encrypted file!\n")
		os.Exit(0)
	}

	return
}
