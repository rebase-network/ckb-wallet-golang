package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
	newblake2b "github.com/liushooter/blake2b"
	"github.com/liushooter/ckb-wallet-golang/ecc"
)

var (
	putf  = fmt.Printf
	putln = fmt.Println
)

const (
	PREFIX_MAINNET string = "ckb"
	PREFIX_TESTNET string = "ckt"
)

func main() {
	seed := rand.Reader

	keyPair, err := ecc.GenerateKey(seed)
	if err != nil {
		panic(err)
	}

	rawPubKey := keyPair.PublicKey

	privBytes := keyPair.ToBytes()
	privKey := byteString(privBytes)

	compressionPubKey := rawPubKey.ToBytes()
	pubKey := byteString(compressionPubKey)

	putf("Privkey: 0x%s\n", privKey)
	putf("Pubkey: 0x%s\n", pubKey)

	blake160 := genBlake160(pubKey)
	putf("Blake160: 0x%x\n", blake160)

	testaddr := genCkbAddr(blake160, PREFIX_TESTNET)
	putf("testAddr: %s\n", testaddr)

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
