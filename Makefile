flags=-X main.buildDate=`date -u '+%Y-%m-%d'` -X main.gitHash=`git rev-parse --short HEAD`

build:
	go build -ldflags "$(flags)" -o ckb-wallet wallet.go

build-all:
	GOARCH=amd64 GOOS=darwin  go build -ldflags "$(flags) -w -s" -o ckb-wallet-mac wallet.go
	GOARCH=amd64 GOOS=linux   go build -ldflags "$(flags) -w -s" -o ckb-wallet-linux wallet.go
	GOARCH=amd64 GOOS=windows go build -ldflags "$(flags) -w -s" -o ckb-wallet-win.exe wallet.go

shasum:
	shasum -a 256 ckb-wallet-linux ckb-wallet-mac ckb-wallet-win.exe