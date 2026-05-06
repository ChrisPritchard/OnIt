go env -w GOOS=linux GOARCH=arm64 CGO_ENABLED=0
go build -o op ./cmd/onit-disp
go env -u GOOS GOARCH CGO_ENABLED
ssh pizero2w "ps aux | grep ./op$ | cut -d' ' -f8 | xargs kill"
scp .\op pizero2w:
ssh pizero2w "chmod +x op && ./op" &
del op