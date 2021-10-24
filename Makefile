GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

export GO111MODULE=on

build:
	$(GOGET)
	$(GOBUILD) homekit-dht22.go

run:
	$(GORUN) homekit-dht22.go
