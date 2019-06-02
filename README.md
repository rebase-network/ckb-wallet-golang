# Ckb Wallet
gen ckb address


# Go Build

```
GOARCH=amd64 GOOS=darwin  go build -ldflags "-w -s" -o ckb-wallet-mac wallet.go
GOARCH=amd64 GOOS=linux   go build -ldflags "-w -s" -o ckb-wallet-linux wallet.go
GOARCH=amd64 GOOS=windows go build -ldflags "-w -s" -o ckb-wallet-win.exe wallet.go
```

# ChangeLog

## v0.2.1
- Support windows10 system
- Add Mainnet address
-  output the project version number