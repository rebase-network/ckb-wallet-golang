package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/btcsuite/btcutil/bech32"
	newblake2b "github.com/liushooter/blake2b"
	"github.com/liushooter/ckb-wallet-golang/ecc"
)

var (
	putf      = fmt.Printf
	putln     = fmt.Println
	gitHash   = ""
	buildDate = ""
)

const (
	VERSION        string = "v0.4.1"
	PREFIX_MAINNET string = "ckb"
	PREFIX_TESTNET string = "ckt"
)

type CkbConfig struct {
	DataDir string `toml:"data_dir"`

	Chain struct {
		Spec string `toml:"spec"`
	} `toml:"chain"`

	Logger struct {
		Filter      string `toml:"filter"`
		Color       bool   `toml:"color"`
		LogToFile   bool   `toml:"log_to_file"`
		LogToStdout bool   `toml:"log_to_stdout"`
	} `toml:"logger"`

	Sentry struct {
		Dsn string `toml:"dsn"`
	} `toml:"sentry"`

	Network struct {
		ListenAddresses             []string      `toml:"listen_addresses"`
		PublicAddresses             []interface{} `toml:"public_addresses"`
		Bootnodes                   []string      `toml:"bootnodes"`
		ReservedPeers               []interface{} `toml:"reserved_peers"`
		ReservedOnly                bool          `toml:"reserved_only"`
		MaxPeers                    int           `toml:"max_peers"`
		MaxOutboundPeers            int           `toml:"max_outbound_peers"`
		PingIntervalSecs            int           `toml:"ping_interval_secs"`
		PingTimeoutSecs             int           `toml:"ping_timeout_secs"`
		ConnectOutboundIntervalSecs int           `toml:"connect_outbound_interval_secs"`
	} `toml:"network"`

	RPC struct {
		ListenAddress      string   `toml:"listen_address"`
		MaxRequestBodySize int      `toml:"max_request_body_size"`
		Modules            []string `toml:"modules"`
	} `toml:"rpc"`

	Sync struct {
		OrphanBlockLimit int `toml:"orphan_block_limit"`
	} `toml:"sync"`

	TxPool struct {
		MaxMemSize          int   `toml:"max_mem_size"`
		MaxCycles           int64 `toml:"max_cycles"`
		MaxVerfifyCacheSize int   `toml:"max_verfify_cache_size"`
	} `toml:"tx_pool"`

	Script struct {
		Runner string `toml:"runner"`
	} `toml:"script"`

	Store struct {
		HeaderCacheSize     int `toml:"header_cache_size"`
		CellOutputCacheSize int `toml:"cell_output_cache_size"`
	} `toml:"store"`

	BlockAssembler struct {
		CodeHash string   `toml:"code_hash"`
		Args     []string `toml:"args"`
	} `toml:"block_assembler"`
}

type Wallet struct {
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
	codehash := flag.String("codehash", "0xf1951123466e4479842387a66fabfd6b65fc87fd84ae8e6cd3053edb27fff2fd", "output codehash")
	config := flag.String("config", "ckb.toml", "output miner config file")

	loop := flag.Bool("loop", false, "")
	num := flag.Int("num", 1, "loop num times")

	flag.Parse()

	if *ver {
		putf("Ckb Wallet Version: %s (%s %s)\n\n", VERSION, gitHash, buildDate)
		os.Exit(0)
	}

	if *loop {
		if *num <= 0 || *num > 1001 {
			putf("-num must 0 ≤ num ≤ 1000\n")
			os.Exit(1001)
		}

		if *privkeyFlag != "" && *num != 1 {
			putf("-privkey -num mutual\n")
			os.Exit(1002)
		}

		for i := 0; i < *num; i++ {
			seed := crand.Reader
			keyPair, err := ecc.GenerateKey(seed)
			if err != nil {
				panic(err)
			}

			rawPubKey := keyPair.PublicKey

			privBytes := keyPair.ToBytes()
			privKey := byteString(privBytes)

			compressionPubKey := rawPubKey.ToBytes()
			pubKey := byteString(compressionPubKey)

			blake160 := genBlake160(pubKey)

			testaddr := genCkbAddr(blake160, PREFIX_TESTNET)
			mainnetaddr := genCkbAddr(blake160, PREFIX_MAINNET)

			putf("0x%s,0x%s,0x%x,%s,%s\n", privKey, pubKey, blake160, testaddr, mainnetaddr)

		}

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

	wallet := Wallet{
		Privkey:     fmt.Sprintf("0x%s", privKey),
		Pubkey:      fmt.Sprintf("0x%s", pubKey),
		Blake160:    fmt.Sprintf("0x%x", blake160),
		TestnetAddr: testaddr,
		MainnetAddr: mainnetaddr,
	}

	if *format == "json" {
		file, _ := json.MarshalIndent(wallet, "", "")
		putf("%s\n", file)
	} else if *format == "csv" {
		putf("0x%s,0x%s,0x%x,%s,%s\n", privKey, pubKey, blake160, testaddr, mainnetaddr)
	} else {
		putf("Privkey: 0x%s\nPubkey: 0x%s\n", privKey, pubKey)
		putf("Blake160: 0x%x\n", blake160)
		putf("TestnetAddr: %s\nMainnetAddr: %s\n", testaddr, mainnetaddr)
	}

	if *config != "" {
		var ckbcfg CkbConfig
		_, err := toml.DecodeFile(*config, &ckbcfg)
		if err != nil {
			panic(err)
		}

		ckbcfg.BlockAssembler.CodeHash = *codehash //
		ckbcfg.BlockAssembler.Args = []string{fmt.Sprintf("0x%x", blake160)}
		ckbfile, _ := ckbcfg.toTOML()

		_ = ioutil.WriteFile("newckb.toml", ckbfile.Bytes(), 0644)

		putf("\nGenerate the miner config file newckb.toml\n")

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

func (t *CkbConfig) toTOML() (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	encoder := toml.NewEncoder(b)

	if err := encoder.Encode(t); err != nil {
		return nil, err
	}
	return b, nil
}
