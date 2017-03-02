.PHONY: build

build:
	rm -rf build/*
	GOOS=linux GOARCH=amd64 go build -o build/polymerase github.com/dollarshaveclub/polymerase
	tar -c -C build polymerase | gzip -c > build/polymerase_linux_amd64.tar.gz