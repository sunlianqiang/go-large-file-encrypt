package rsacrypt

import (
	"bufio"
	"log"
	"os"

	"github.com/astaxie/beego/logs"
)

func GetPublickKeys(pubKeyFile string) (pubKeysBytes map[int]string, err error) {
	pubKeysBytes = make(map[int]string)

	var readFile *os.File
	readFile, err = os.Open(pubKeyFile)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	count := 0
	for fileScanner.Scan() {
		pubKeysBytes[count] = fileScanner.Text()
		count++
		// pubKeysBytes = append(pubKeysBytes, fileScanner.Text())
	}

	readFile.Close()

	for k, eachline := range pubKeysBytes {
		logs.Debug("public key %v: %v\n", k, eachline)
	}

	return
}
