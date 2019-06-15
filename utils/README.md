```
GOARCH=amd64 GOOS=darwin  go build -ldflags "-w -s" -o str2hex-mac str2hex.go
GOARCH=amd64 GOOS=linux   go build -ldflags "-w -s" -o str2hex-linux str2hex.go
GOARCH=amd64 GOOS=windows go build -ldflags "-w -s" -o str2hex-win.exe str2hex.go
```