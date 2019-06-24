# Ckb Wallet
gen ckb address

## Wallet

```
./ckb-wallet-mac -v

./ckb-wallet-mac -format txt

./ckb-wallet-mac -privkey 1234

./ckb-wallet-mac -codehash 0x5678

./ckb-wallet-mac -privkey 1234 -codehash 0x5678 -format csv

./ckb-wallet-mac -config ckb.toml -privkey 1234 -codehash 0x5678 -format json

```

`./ckb-wallet-mac -format txt | grep Blake160 | grep -oE "[^:]+$"`


## Go Build

```
flags="-X main.buildDate=`date -u '+%Y-%m-%d'` -X main.gitHash=`git rev-parse --short HEAD`"

GOARCH=amd64 GOOS=darwin  go build -ldflags "$flags -w -s" -o ckb-wallet-mac wallet.go
GOARCH=amd64 GOOS=linux   go build -ldflags "$flags -w -s" -o ckb-wallet-linux wallet.go
GOARCH=amd64 GOOS=windows go build -ldflags "$flags -w -s" -o ckb-wallet-win.exe wallet.go
```

## Verifying the Release

```
shasum -a 256 ckb-wallet-linux ckb-wallet-mac ckb-wallet-win.exe
```

## ChangeLog

### v0.4.1
- update version display

### v0.4
- generate newckb.toml

### v0.3.3
- update codehash

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
