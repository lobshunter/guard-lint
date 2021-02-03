GOENV   := GO111MODULE=on CGO_ENABLED=0
GO 		:= $(GOENV) go

LDFLAGS := -w -s

YAML_PATH := examples

default: vet lint-schema unique-name

vet:
	$(GO) vet ./...

lint-schema:
	$(GO) build -ldflags '$(LDFLAGS)' -o  output/lint lint-tools/lint-schema/main.go
unique-name:
	$(GO) build -ldflags '$(LDFLAGS)' -o  output/unique lint-tools/unique-name/main.go

lint:
	@# git diff-tree will get added/modified files in last commit
	@# grep will extract files under $(YAML_PATH) folder and ends with yaml/yml (i.e. only lint yaml files under example folder)
	@output/lint schema.json `git diff-tree --name-only --no-commit-id --diff-filter=AM -r HEAD | grep -E "^$(YAML_PATH).+(yaml|yml)$$"` 

check-unique:
	@output/unique `find $(YAML_PATH) -name '*.yaml' -o -name '*.yml'`