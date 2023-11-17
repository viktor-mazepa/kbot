APP=$(shell basename -s .git $(shell git remote get-url origin))
REGISTRY=europe-central2-docker.pkg.dev/devops2023-405122/kbot-repo
VERSION=${shell git describe --tags --abbrev=0}-${shell git rev-parse --short HEAD}
TARGETOS=linux
TARGETARCH=amd64

format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v 
get:
	go get

#default build
build: format get
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v -o kbot -ldflags "-X 'github.com/viktor-mazepa/kbot/cmd.appVersion=${VERSION}'" 

windows: format get
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -o kbot -ldflags "-X 'github.com/viktor-mazepa/kbot/cmd.appVersion=${VERSION}'" 

windows386: format get
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -v -o kbot -ldflags "-X 'github.com/viktor-mazepa/kbot/cmd.appVersion=${VERSION}'" 

linux: format get
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o kbot -ldflags "-X 'github.com/viktor-mazepa/kbot/cmd.appVersion=${VERSION}'" 

mac: format get
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -o kbot -ldflags "-X 'github.com/viktor-mazepa/kbot/cmd.appVersion=${VERSION}'" 

image:
	docker build . -t $(REGISTRY)/$(APP):$(VERSION)-$(TARGETARCH) --build-arg TARGETARCH=$(TARGETARCH)

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

clean: 
	rm -rf kbot
	docker rmi ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}



