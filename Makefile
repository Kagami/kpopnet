all: install-deps

env:
	virtualenv -p python3 --system-site-packages env

install-deps: env
	env/bin/pip install -e .[tests]

lint:
	env/bin/flake8 kpopnet
