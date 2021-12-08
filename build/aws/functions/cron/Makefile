all:
	mkdir -p build
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o build/main ../../../../aws/functions/cron/main.go
	cd build && zip function.zip main

clean:
	$(RM) -rf build