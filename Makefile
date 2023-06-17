cert:
	mkdir -p cert && \
	openssl req -x509 -nodes -newkey rsa:2048 -keyout cert/key.pem -out cert/cert.pem -subj "//CN=localhost" -days 365

test:
	go test -cover -coverprofile cover.out ./...

build:
	docker build -t microbin:dev .

run:
	docker run --rm --name microbin -p 8080:8080 microbin:dev

runtls:
	MSYS_NO_PATHCONV=1 docker run --rm --name microbin -p 8080:8080 \
	-v ${PWD}/cert:/var/cert:ro \
	-e TLS_ENABLED=true \
	-e TLS_KEY_FILE=/var/cert/key.pem \
	-e TLS_CERT_FILE=/var/cert/cert.pem \
	-e LOG_FORMAT_JSON=false \
	microbin:dev

.PHONY: cert test build run runtls