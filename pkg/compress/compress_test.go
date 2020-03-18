package compress

import (
	"fmt"
	"testing"
	"time"
)

func TestZip(t *testing.T) {

	// List of Files to Zip
	infile := "file1.txt"
	midfile := "/Users/slq/Downloads/Microsoft_Remote_Desktop_Beta.app.zip"
	// bigfile := "/Users/slq/Downloads/Apollo.11.2019.1080p.BluRay.x264.DTS-HD.MA.5.1-FGT/Apollo.11.2019.1080p.BluRay.x264.DTS-HD.MA.5.1-FGT.mkv"
	files := []string{midfile, infile, "file2.txt"}
	zipFile := infile + "." + time.Now().Local().Format(time.RFC3339) + ".zip"

	if err := ZipFiles(zipFile, files); err != nil {
		panic(err)
	}
	fmt.Println("Zipped File:", zipFile)
}
