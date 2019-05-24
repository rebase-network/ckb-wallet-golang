package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
	newblake2b "github.com/liushooter/blake2b"
)

const (
	PREFIX_MAINNET string = "ckb"
	PREFIX_TESTNET string = "ckt"
)

func main() {

	curve := elliptic.P256() // http://golang.org/pkg/crypto/elliptic/#P256

	keyPair, err := ecdsa.GenerateKey(curve, crand.Reader) // this generates a public & private key pair
	if err != nil {
		panic(err)
	}

	pubKey := keyPair.PublicKey

	fmt.Printf("PrivateKey : 0x%x\n", keyPair.D)

	var compressionPubKey string
	if pubKey.Y.Bit(0) == 0 { //even
		compressionPubKey = fmt.Sprintf("%s%x", "02", pubKey.X)
	} else { // odd
		compressionPubKey = fmt.Sprintf("%s%x", "03", pubKey.X)
	}

	fmt.Printf("Pubkey: 0x%s\n", compressionPubKey)

	hexbin, _ := hex.DecodeString(compressionPubKey)

	ckbsum := newblake2b.CkbSum256(hexbin)
	blake160 := ckbsum[:20]

	fmt.Printf("Blake160: 0x%x\n", blake160)

	testaddr := genCkbAddr(blake160, PREFIX_TESTNET)
	fmt.Printf("testaddr: %s\n", testaddr)

}

func genCkbAddr(blake160Addr []byte, prefix string) string {

	typebin, _ := hex.DecodeString("01")
	bin_idx := []byte("P2PH")

	payload := append(typebin, bin_idx...)
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

// https://www.socketloop.com/tutorials/golang-example-for-ecdsa-elliptic-curve-digital-signature-algorithm-functions
