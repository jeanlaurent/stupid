depends:
	dep ensure

test:
	go test -v .

release:
	rm -rf build
	mkdir build
	GOOS=darwin GOARCH=amd64 go build -o build/stupid
	GOOS=windows GOARCH=amd64 go build -o build/stupid.exe