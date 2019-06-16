build:
	go build -ldflags "-w -s" -o ckb-wallet wallet.go

build-all:
	GOARCH=amd64 GOOS=darwin  go build -ldflags "-w -s" -o ckb-wallet-mac wallet.go
	GOARCH=amd64 GOOS=linux   go build -ldflags "-w -s" -o ckb-wallet-linux wallet.go
	GOARCH=amd64 GOOS=windows go build -ldflags "-w -s" -o ckb-wallet-win.exe wallet.go

shasum:
	shasum -a 256 ckb-wallet-linux ckb-wallet-mac ckb-wallet-win.exe