# Microbin

[![CI Build](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/ci-build.yaml)

Yet another httpbin

**Deploy**

    docker run --rm -d --name microbin -p8080:8080 ghcr.io/sshilin/microbin:latest

or

    kubectl apply -f microbin.yaml

**API**

- GET /

Output:
```
microbin v.dev
```

- GET /headers

Output:
```
{
  "pod": {
    "name": "microbin-6d66c6d744-vkzcb",
    "namespace": "default",
    "node": "docker-desktop"
  },
  "headers": {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
    "Accept-Encoding": "gzip, deflate",
    "Accept-Language": "en-US,en;q=0.9,ru;q=0.8",
    "Connection": "keep-alive",
    "Dnt": "1",
    "Upgrade-Insecure-Requests": "1",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
  }
}
```