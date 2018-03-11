export GOPATH = $(PWD)/go

all: pydeps go/bin/kpopnet

py/env:
	virtualenv -p python3 --system-site-packages $@

pydeps: py/env
	$</bin/pip install -e .[tests]

pylint: py/env
	$</bin/flake8 py/kpopnet

tswatch:
	npm start

go/bin/kpopnet: go/src/kpopnet/**/*
	go get -v kpopnet

goserve: go/bin/kpopnet
	$<

gofmt:
	go fmt kpopnet/...
