default:
	echo "make what? cli, cron, api or all?"
	exit 2

all:
	$(MAKE) cli
	$(MAKE) cron
	$(MAKE) api

cli:
	mkdir -p build/cli
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/cli/linux-amd64 ../../aws/functions/cli/main.go
	env GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o build/cli/linux-arm64 ../../aws/functions/cli/main.go
	env GOOS=darwin GOARCH=amd64 go build -o build/cli/darwin-amd64 ../../aws/functions/cli/main.go
	env GOOS=darwin GOARCH=arm64 go build -o build/cli/darwin-arm64 ../../aws/functions/cli/main.go
	env GOOS=windows GOARCH=amd64 go build -o build/cli/windows-amd64 ../../aws/functions/cli/main.go

cron:
	mkdir -p build/cron
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/cron/main ../../aws/functions/cron/main.go
	cd build/cron && zip function.zip main

api:
	mkdir -p build/api
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/api/main ../../aws/functions/api/main.go
	cd build/api && zip function.zip main

clean:
	$(RM) -rf build