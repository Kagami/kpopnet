export GOPATH = $(PWD)/go
export PATH := $(PATH):$(PWD)/go/bin

all: pydeps tsdeps go/bin/kpopnetd
test: pylint tslint gofmt-staged go/bin/kpopnetd gotest

py/env:
	virtualenv -p python3 --system-site-packages $@

pydeps: py/env
	$</bin/pip install -e .[tests]

pylint: py/env
	$</bin/flake8 py/kpopnet

tsdeps:
	npm install

tswatch:
	npm start

tslint:
	npm -s test

go/bin/go-bindata:
	go get github.com/kevinburke/go-bindata/...

GOSRC = $(shell find go/src/kpopnet -type f)
go/bin/kpopnetd: go/bin/go-bindata $(GOSRC)
	go generate kpopnet/...
	go get -v kpopnet/...

goserve: go/bin/kpopnetd
	$< serve

gofmt:
	go fmt kpopnet/...

gofmt-staged:
	./gofmt-staged.sh

gotags:
	ctags -R go/src/kpopnet

go/src/kpopnet/testdata:
	git clone https://github.com/Kagami/go-face-testdata go/src/kpopnet/testdata

gotest: go/src/kpopnet/testdata
	go test -v kpopnet
