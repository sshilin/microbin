# microbin

[![CI Build](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml)

Httpbin-like simple application with intention to be deployed on the Kubernetes cluster, adds pod info to json responses.

Goals:
- investigate traffic routing in K8s cluster
- play with GitHub Actions and Packages

### Deployment
---

    docker run --rm -d --name microbin -p8080:8080 ghcr.io/sshilin/microbin:latest

or

    kubectl apply -f microbin.yaml

### Usage
---
- GET / - prints version
```
microbin v.dev
```
- GET /headers - prints request headers and pod info
```
{
  "pod": {
    "name": "microbin-6d66c6d744-vkzcb",
    "namespace": "default",
    "node": "docker-desktop"
  },
  "headers": {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*",
    "Accept-Encoding": "gzip, deflate",
    "Accept-Language": "en-US,en;q=0.9,ru;q=0.8",
    "Connection": "keep-alive",
    "Dnt": "1",
    "Upgrade-Insecure-Requests": "1",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"
  }
}
```
