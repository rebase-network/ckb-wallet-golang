# Ckb Wallet
gen ckb address

`./ckb-wallet-mac -format txt | grep Blake160 | grep -oE "[^:]+$"`

# Go Build

```
GOARCH=amd64 GOOS=darwin  go build -ldflags "-w -s" -o ckb-wallet-mac wallet.go
GOARCH=amd64 GOOS=linux   go build -ldflags "-w -s" -o ckb-wallet-linux wallet.go
GOARCH=amd64 GOOS=windows go build -ldflags "-w -s" -o ckb-wallet-win.exe wallet.go
```

# Verifying the Release

```
shasum -a 256 ckb-wallet-linux
shasum -a 256 ckb-wallet-mac
shasum -a 256 ckb-wallet-win.exe
```

# ChangeLog

## v0.3.2

## v0.3.1
- output `format` support csv

## v0.3
- use golang pkg `flag` refactoring project
- add subcommand `-v` `-format` `-privkey` `-config`
- Generate the miner config file ckb.toml

## v0.2.1
- Support windows10 system
- Add Mainnet address
- output the project version number
