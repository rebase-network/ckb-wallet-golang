package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	newblake2b "github.com/liushooter/blake2b"
)

func main() {

	pubkeyCurve := elliptic.P256() // http://golang.org/pkg/crypto/elliptic/#P256

	keyPair, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
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

	blake160 := newblake2b.CkbSum256(hexbin)

	fmt.Printf("Blake160: 0x%x\n", blake160[:20])
}

// https://www.socketloop.com/tutorials/golang-example-for-ecdsa-elliptic-curve-digital-signature-algorithm-functions
