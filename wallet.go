package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"

	"github.com/btcsuite/btcutil/bech32"
	newblake2b "github.com/liushooter/blake2b"
	"github.com/liushooter/ckb-wallet-golang/ecc"
)

var (
	putf  = fmt.Printf
	putln = fmt.Println
)

const (
	version        string = "v0.3"
	PREFIX_MAINNET string = "ckb"
	PREFIX_TESTNET string = "ckt"
)

type output struct {
	Privkey     string
	Pubkey      string
	Blake160    string
	TestnetAddr string
	MainnetAddr string
}

func main() {
	ver := flag.Bool("v", false, "show version and exit")
	privkeyFlag := flag.String("privkey", "", "ehter privkey")
	format := flag.String("format", "json", "output format")
	config := flag.String("config", "0x9e3b3557f11b2b3532ce352bfe8017e9fd11d154c4c7f9b7aaaa1e621b539a08", "output miner config file")

	flag.Parse()

	if *ver {
		putf("Ckb Wallet Version: %s\n\n", version)
		os.Exit(0)
	}

	var keyPair ecc.PrivateKey

	if *privkeyFlag != "" {
		bignum := new(big.Int)
		bignum.SetString(*privkeyFlag, 16)
		keyPair = *ecc.NewPrivateKey(bignum)
	} else {
		var err error
		seed := crand.Reader
		keyPair, err = ecc.GenerateKey(seed)
		if err != nil {
			panic(err)
		}
	}

	rawPubKey := keyPair.PublicKey

	privBytes := keyPair.ToBytes()
	privKey := byteString(privBytes)

	compressionPubKey := rawPubKey.ToBytes()
	pubKey := byteString(compressionPubKey)

	blake160 := genBlake160(pubKey)

	testaddr := genCkbAddr(blake160, PREFIX_TESTNET)
	mainnetaddr := genCkbAddr(blake160, PREFIX_MAINNET)

	if *format == "json" {
		data := output{
			Privkey:     fmt.Sprintf("0x%s", privKey),
			Pubkey:      fmt.Sprintf("0x%s", pubKey),
			Blake160:    fmt.Sprintf("0x%x", blake160),
			TestnetAddr: testaddr,
			MainnetAddr: mainnetaddr,
		}

		file, _ := json.MarshalIndent(data, "", "")
		fmt.Printf("%s\n", file)

	} else {
		putf("Privkey: 0x%s\nPubkey: 0x%s\n", privKey, pubKey)
		putf("Blake160: 0x%x\n", blake160)
		putf("TestAddr: %s\nMainnetAddr: %s\n\n", testaddr, mainnetaddr)
	}

	if *config != "" {
		ckbfile := fmt.Sprintf("\n[block_assembler]\ncode_hash = '%s'\nargs = ['%s']\n", *config, fmt.Sprintf("0x%x", blake160))
		putf("\nAlready Generate the miner config file ckb.toml\n")
		putf(ckbfile)

		_ = ioutil.WriteFile("ckb.toml", []byte(ckbfile), 0644)
	}

	if runtime.GOOS == "windows" {
		fmt.Scanln() // Enter Key to terminate the console screen
	}
}

func genBlake160(pubKey string) []byte {
	hexbin, _ := hex.DecodeString(pubKey)

	ckbsum := newblake2b.CkbSum256(hexbin)
	blake160 := ckbsum[:20]
	return blake160
}

func genCkbAddr(blake160Addr []byte, prefix string) string {

	typebin, _ := hex.DecodeString("01")
	flag := []byte("P2PH")

	payload := append(typebin, flag...)
	payload = append(payload, blake160Addr...)

	converted, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		panic(err)
	}

	addr, err := bech32.Encode(prefix, converted)
	if err != nil {
		panic(err)
	}

	return addr
}

func byteString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02x", b[i])
	}
	return s
}
