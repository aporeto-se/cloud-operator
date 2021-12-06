default:
	echo "make what? cli, all"
	exit 2

all:
	$(MAKE) cli

cli:
	mkdir -p build/cli
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/cli/linux-amd64 ../../gcp/functions/cli/main.go
	env GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o build/cli/linux-arm64 ../../gcp/functions/cli/main.go
	env GOOS=darwin GOARCH=amd64 go build -o build/cli/darwin-amd64 ../../gcp/functions/cli/main.go
	env GOOS=darwin GOARCH=arm64 go build -o build/cli/darwin-arm64 ../../gcp/functions/cli/main.go
	env GOOS=windows GOARCH=amd64 go build -o build/cli/windows-amd64 ../../gcp/functions/cli/main.go

clean:
	$(RM) -rf build
