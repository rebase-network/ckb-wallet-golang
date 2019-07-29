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

	"github.com/rebase-network/ckb-wallet-golang/ecc"

	"github.com/BurntSushi/toml"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/gookit/color"
	newblake2b "github.com/liushooter/blake2b"
)

var (
	putf  = fmt.Printf
	putln = fmt.Println
	puts  = fmt.Sprintf

	gitHash   = ""
	buildDate = ""

	red   = color.FgRed.Render
	green = color.FgGreen.Render
)

const (
	VERSION        = "v0.5.1"
	PREFIX_MAINNET = "ckb"
	PREFIX_TESTNET = "ckt"

	ErrColor = "\033[1;31m%s\033[0m"
)

type CkbConfig struct {
	// https://xuri.me/toml-to-go/

	DataDir string `toml:"data_dir"`

	Chain struct {
		Spec struct {
			Bundled string `toml:"bundled"`
		} `toml:"spec"`
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
		Upnp                        bool          `toml:"upnp"`
		DiscoveryLocalAddress       bool          `toml:"discovery_local_address"`
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
		Data     string   `toml:"data"`
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
	codehash := flag.String("codehash", "0x94334bdda40b69bae067d84937aa6bbccf8acd0df6626d4b9ac70d4612a11933", "output codehash")
	config := flag.String("config", "ckb.toml", "output config file")
	data := flag.String("data", "", "set ckb cellbase data")

	loop := flag.Bool("loop", false, "")
	num := flag.Int("num", 1, "loop num times")

	flag.Parse()

	if *ver {
		putf("Ckb Wallet Version: %s (%s %s)\n\n", VERSION, gitHash, buildDate)
		os.Exit(0)
	}

	if *loop {
		if *num <= 0 || *num > 1001 {
			putf(ErrColor, "-num must 0 ≤ num ≤ 1000\n")
			os.Exit(1001)
		}

		if *privkeyFlag != "" && *num != 1 {
			putf(ErrColor, "-privkey -num mutual\n")
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

			blake160 := genBlake160(compressionPubKey)

			testaddr := genCkbAddr(PREFIX_TESTNET, blake160)
			mainnetaddr := genCkbAddr(PREFIX_MAINNET, blake160)

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

	blake160 := genBlake160(compressionPubKey)

	testaddr := genCkbAddr(PREFIX_TESTNET, blake160)
	mainnetaddr := genCkbAddr(PREFIX_MAINNET, blake160)

	wallet := Wallet{
		Privkey:     puts("0x%s", privKey),
		Pubkey:      puts("0x%s", pubKey),
		Blake160:    puts("0x%x", blake160),
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

		if !isExists(*config) {
			putf(red("\n", *config, " File not exists\n"))
			os.Exit(1003)
		}

		var ckbcfg CkbConfig
		_, err := toml.DecodeFile(*config, &ckbcfg)
		if err != nil {
			panic(err)
		}

		ckbcfg.BlockAssembler.CodeHash = *codehash // codehash
		ckbcfg.BlockAssembler.Args = []string{fmt.Sprintf("0x%x", blake160)}
		ckbcfg.BlockAssembler.Data = *data // cellbase data

		ckbfile, _ := ckbcfg.toTOML()

		_ = ioutil.WriteFile("newckb.toml", ckbfile.Bytes(), 0644)

		putf(green("\nGenerate the config file: newckb.toml\n"))

	}

	if runtime.GOOS == "windows" {
		fmt.Scanln() // Enter Key to terminate the console screen
	}
}

func genBlake160(pubKeyBin []byte) []byte {

	ckbsum := newblake2b.CkbSum256(pubKeyBin)
	blake160 := ckbsum[:20]
	return blake160
}

func genCkbAddr(prefix string, blake160Addr []byte) string {

	typebin, _ := hex.DecodeString("01")
	flag, _ := hex.DecodeString("00")

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

func isExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
