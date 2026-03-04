APP=miroxy
VERSION?=0.1.0

.PHONY: build build-all clean

build:
	go build -ldflags="-s -w" -o $(APP) .

build-pi2:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o dist/$(APP)-linux-arm .

# Pi Zero 2 W = ARM v6/v7, Pi 3/4 = ARM64
build-all:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o dist/$(APP)-linux-arm .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(APP)-linux-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(APP)-linux-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(APP)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(APP)-darwin-amd64 .

clean:
	rm -rf $(APP) dist/
