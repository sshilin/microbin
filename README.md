# Microbin

Yet another httpbin

**Multi-Arch build**

    docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKERHUB_USER/microbin:latest --build-arg BUILD_VERSION=v0.0.0-$(git rev-parse --short HEAD) --push .

**Run**

    docker run --rm -d --name microbin -p8080:8080 $DOCKERHUB_USER/microbin:latest

