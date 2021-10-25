.DEFAULT_GOAL:=build

BINARY:=server

clean:
	rm -f ${BINARY}
	rm -rf docs
	rm -rf dist

prep:
	mkdir dist/linux
	mkdir dist/darwin
	mkdir dist/windows

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o dist/linux/${BINARY}

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o dist/darwin/${BINARY}

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows/${BINARY}.exe

build: clean docs build-linux build-darwin build-windows

deps:
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/spf13/cobra/cobra

docs:
	swag init -g cmd/root.go

test: build
	./${BINARY}

