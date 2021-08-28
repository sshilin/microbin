[![CI Build](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml)

Microbin is a simple containerazed http(s) server which replys with the request headers and additional info about the application instance.

### Usage
---
Endpoints:
- `GET /` - return version
- `GET /headers` - return request headers and pod info in JSON
```
$ curl http://localhost:8080/headers
{
  "proto": "HTTP/1.1",
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/7.77.0"
  },
  "pod": {
    "name": "microbin-58cf7dff8d-7ts8l",
    "namespace": "default",
    "node": "docker-desktop"
  }
}
```
### Environment Variables
---
The table below provides an overview of optional environment variables that can be used to configure microbin.

| Key                 |  Description                | Default         |
|:--------------------|:----------------------------|:----------------|
| `LISTEN`            | Listen on host:port         | 0.0.0.0:8080    |
| `TLS_ENABLED`       | Enable TLS                  | false           |
| `TLS_KEY_FILE`      | TLS key filepath            | ""              |
| `TLS_CERT_FILE`     | TLS cert filepath           | ""              |

### Deployment
---
Starting Docker container:

    docker run --rm -d -p 8080:8080 --name microbin ghcr.io/sshilin/microbin:latest

Deploying on Kubernetes:

    kubectl apply -f microbin.yaml

### Enable HTTPS
---
This example shows how to enable HTTPS with a self signed certificate. A self signed cert can be generated via a single openssl command

    openssl req -x509 -nodes -newkey rsa:4096 -keyout ./certs/key.pem -out ./certs/cert.pem -subj "//CN=localhost" -days 365

**Docker**

Run a container with mounted key and cert files. (MSYS_NO_PATHCONV=1 temparary disables path conversion in GitBash/Msys2)

```bash
$ MSYS_NO_PATHCONV=1 \
  docker run --rm -d -p 8080:8080 --name microbin \
  -v "$(pwd)/certs:/var/tls:ro" \
  -e TLS_ENABLED=true \
  -e TLS_KEY_FILE=/var/tls/key.pem \
  -e TLS_CERT_FILE=/var/tls/cert.pem \
  microbin:latest
```

**Kubernetes**

First upload the cert and key to the `microbin-certs` secret, then deploy `microbin-https.yaml`

    kubectl create secret generic microbin-certs --from-file=./certs

    kubectl apply -f microbin-https.yaml
