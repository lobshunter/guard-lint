GOENV   := GO111MODULE=on CGO_ENABLED=0
GO 		:= $(GOENV) go

LDFLAGS := -w -s

default: vet build

vet:
	$(GO) vet ./...

build:
	$(GO) build -ldflags '$(LDFLAGS)' -o  output/lint .