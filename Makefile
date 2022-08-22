test:
	cd app && go test -count 1 ./...

build:
	docker build -t github.com/sshilin/microbin:dev .

.PHONY: test build