default: build

build: fix
	go build -v -ldflags '-s -w' .

static: fix
	go build -v -tags netgo -ldflags '-s -w -linkmode external -extldflags "-static"' .

fix: *.go
	goimports -l -w .
	gofmt -l -w -s .
