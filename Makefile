.DEFAULT_GOAL := .make

prepare:
	@scripts/prepare.sh

lint: prepare
	@scripts/lint.sh

build: prepare
	@scripts/build.sh

.make: clean prepare lint build

clean:
	@scripts/clean.sh
