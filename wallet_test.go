package main

import (
	"math/big"
	"testing"

	"github.com/rebase-network/ckb-wallet-golang/ecc"
	"github.com/stretchr/testify/assert"
)

var (
	// https://github.com/nervosnetwork/ckb-sdk-ruby/blob/develop/spec/ckb/address_spec.rb

	testPrivkey  = "0xe79f3207ea4980b7fed79956d5934249ceac4751a4fae01a0f7c4a96884bc4e3"
	testPubkey   = "0x024a501efd328e062c8675f2365970728c859c592beeefd6be8ead3d901330bc01"
	testBlake160 = "0x36c329ed630d6ce750712a477543672adab57f4c"

	ckbAddr = "ckt1qyqrdsefa43s6m882pcj53m4gdnj4k440axqswmu83"
)

func TestWalletAddr(t *testing.T) {
	assert := assert.New(t)

	bignum := new(big.Int)
	bignum.SetString(testPrivkey[2:], 16)

	keyPair := *ecc.NewPrivateKey(bignum)

	rawPubKey := keyPair.PublicKey

	compressionPubKey := rawPubKey.ToBytes()
	pubKey := byteString(compressionPubKey)

	putf("privkey: %s\n", testPrivkey)

	hexPubKey := puts("0x%s", pubKey)
	putf("pubKey: %s\n", hexPubKey)
	assert.Equal(hexPubKey, testPubkey)

	blake160 := genBlake160(compressionPubKey)

	hexBlake160 := puts("0x%x", blake160)
	putf("blake160: %s\n", hexBlake160)
	assert.Equal(hexBlake160, testBlake160)

	testaddr := genCkbAddr(PREFIX_TESTNET, blake160)
	putf("newTestAddr: %s\n", testaddr)
	assert.Equal(testaddr, ckbAddr)

}
