GOENV   := GO111MODULE=on CGO_ENABLED=0
GO 		:= $(GOENV) go

LDFLAGS := -w -s

default: vet build

vet:
	$(GO) vet ./...

build:
	$(GO) build -ldflags '$(LDFLAGS)' -o  output/lint .

lint:
	@# git diff-tree will get added/modified files in last commit
	@# grep will extract files under examples folder and ends with yaml/yml (only lint yaml files under example folder)
	@git diff-tree --name-only --no-commit-id --diff-filter=AM -r HEAD | grep -E "^examples.+(yaml|yml)$$" | xargs output/lint schema.json
