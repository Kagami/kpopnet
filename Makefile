all: install-deps

py/env:
	virtualenv -p python3 --system-site-packages $@

install-deps: py/env
	$^/bin/pip install -e .[tests]

lint: py/env
	$^/bin/flake8 py/kpopnet
