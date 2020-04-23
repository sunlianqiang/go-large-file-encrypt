package gpg

/**
  Example hack to encrypt a file using a GPG encryption key. Works with GPG v2.x.
  The encrypted file e.g. /tmp/data.txt.gpg can then be decrypted using the standard command
  gpg /tmp/data.txt.gpg
  Assumes you have **created** an encryption key and exported armored version.
  You have to read the armored key directly as Go cannot read pubring.kbx (yet).
  Export your key using command:
    gpg2 --export --armor [KEY ID] > /tmp/pubKey.asc
*/

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"io/ioutil"
	"os/exec"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"github.com/astaxie/beego/logs"
)

// change as required

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	return wc.Close()
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func EncryptPGP(msg2Enc string, pubkeyFile string, outfile string) (err error) {

	log.Printf("Public key file:%v, outfile:%v\n", pubkeyFile, outfile)

	// Read in public key
	recipient, err := readEntity(pubkeyFile)
	if err != nil {
		fmt.Printf("EncryptPGP readEntity err:%v\n", err)
		return
	}

	dst, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("EncryptPGP Create outfile err:%v\n", err)
		return
	}
	defer dst.Close()
	err = encrypt([]*openpgp.Entity{recipient}, nil, strings.NewReader(msg2Enc), dst)
	if err != nil {
		fmt.Printf("EncryptPGP, Create outfile err:%v\n", err)
		return
	}

	return
}

func DecryptGPG(encryptedAesKeyFile string) (aesKeyBytes []byte, err error) {
	cmd := exec.Command("gpg", "-d", encryptedAesKeyFile)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logs.Error(err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.Error(err)
		return
	}
	if err = cmd.Start(); err != nil {
		logs.Error(err)
		return
	}

	errBytes, err := ioutil.ReadAll(stderr)
	if nil != err {
		logs.Error(err)
		return
	}
	aesKeyBytes, err = ioutil.ReadAll(stdout)
	if nil != err {
		logs.Error(err)
		return
	}
	fmt.Printf("DecryptAesKey, stderr:%s\n", errBytes)

	if err = cmd.Wait(); err != nil {
		logs.Error(err)
		return
	}

	return
}

