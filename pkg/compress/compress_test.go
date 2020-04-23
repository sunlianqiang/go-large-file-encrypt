package compress

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestZip(t *testing.T) {

	// List of Files to Zip
	infile := "file1.txt"
	files := []string{infile, "file2.txt"}
	zipFile := infile + "." + time.Now().Local().Format(time.RFC3339) + ".zip"

	zipFile = "test.zip"
	if err := ZipFiles(zipFile, files); err != nil {
		panic(err)
	}
	fmt.Println("Zipped File:", zipFile)
}

func TestUnzip(t *testing.T) {
	files, err := Unzip("test.zip", "test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
}
