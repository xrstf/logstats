default: build

build: fix
	go build -v .

fix: *.go
	goimports -l -w .
	gofmt -l -w -s .

app:
	cd cmd/logstats && make

static:
	cd cmd/logstats && make static

deps:
	go get gopkg.in/yaml.v2
	go get golang.org/x/tools/cmd/goimports
