# Ckb Wallet
gen ckb address

## Wallet

```
./ckb-wallet-mac -v

./ckb-wallet-mac -format txt
./ckb-wallet-mac -format csv

./ckb-wallet-mac -privkey 1234
./ckb-wallet-mac -privkey 1234 -format txt

./ckb-wallet-mac -config 0x5678
./ckb-wallet-mac -config 0x5678 -format csv

./ckb-wallet-mac -privkey 1234 -config 0x5678 -format txt
```

`./ckb-wallet-mac -format txt | grep Blake160 | grep -oE "[^:]+$"`


## Go Build

```
GOARCH=amd64 GOOS=darwin  go build -ldflags "-w -s" -o ckb-wallet-mac wallet.go
GOARCH=amd64 GOOS=linux   go build -ldflags "-w -s" -o ckb-wallet-linux wallet.go
GOARCH=amd64 GOOS=windows go build -ldflags "-w -s" -o ckb-wallet-win.exe wallet.go
```

## Verifying the Release

```
shasum -a 256 ckb-tool-linux
shasum -a 256 ckb-tool-mac
shasum -a 256 ckb-tool-win.exe
```

## ChangeLog


### v0.3.3
- update code_hash

### v0.3.2

### v0.3.1
- output `format` support csv

### v0.3
- use golang pkg `flag` refactoring project
- add subcommand `-v` `-format` `-privkey` `-config`
- Generate the miner config file ckb.toml

### v0.2.1
- Support windows10 system
- Add Mainnet address
- output the project version number
