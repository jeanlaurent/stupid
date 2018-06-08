build: test cross

depends:
	dep ensure

test:
	go test -coverprofile=coverage.txt -covermode=atomic -v ./...

bin/%: cmd/%
	go build -o $@ ./$<

bin/stupid-linux: cmd/stupid
	GOOS=linux go build -o $@ ./$<

bin/stupid-darwin: cmd/stupid
	GOOS=darwin go build -o $@ ./$<

bin/stupid-windows.exe: cmd/stupid
	GOOS=windows go build -o $@ ./$<

cross: bin/stupid-linux bin/stupid-darwin bin/stupid-windows.exe
