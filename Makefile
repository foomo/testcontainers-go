.DEFAULT_GOAL:=help
-include .makerc

# --- Config -----------------------------------------------------------------

GOMODS=$(shell find . -type f -name go.mod)
# Newline hack for error output
define br


endef

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: go.work .mise .lefthook
	@:

# Ensure file
go.work:
	@echo "〉initializing go work"
	@go work init
	@go work use -r .
	@go work sync

.PHONY: .mise
# Install dependencies
.mise:
ifeq (, $(shell command -v mise))
	$(error $(br)$(br)Please ensure you have 'mise' installed and activated!$(br)$(br)  $$ brew update$(br)  $$ brew install mise$(br)$(br)See the documentation: https://mise.jdx.dev/getting-started.html)
endif
	@mise install

.PHONY: .lefthook
# Configure git hooks for lefthook
.lefthook:
	@lefthook install --reset-hooks-path

### Tasks

.PHONY: check
## Run lint & tests
check: tidy generate lint.fix test audit

.PHONY: lint
## Run linter
lint:
	@echo "〉go lint"
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && golangci-lint run) &&) true

.PHONY: lint.fix
## Run linter and fix violations
lint.fix:
	@echo "〉go lint with --fix"
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && golangci-lint run --fix) &&) true

.PHONY: generate
## Run go generate
generate:
	@echo "〉go generate"
	@go generate work

.PHONY: test
## Run tests
test:
	@echo "〉go test"
	@GO_TEST_TAGS=-skip go test -tags=safe -coverprofile=coverage.out work

.PHONY: test.race
## Run tests with -race
test.race:
	@GO_TEST_TAGS=-skip go test -tags=safe -coverprofile=coverage.out -race work

### Security

.PHONY: audit
## Run security audit
audit:
	@echo "〉security audit"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && govulncheck ./...) &&) true

### Dependencies

.PHONY: tidy
## Run go mod tidy
tidy:
	@echo "〉go mod tidy"
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && go mod tidy) &&) true
	@go work sync

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@echo "〉go mod outdated"
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && go list -u -m -json all | go-mod-outdated -update -direct) &&) true

.PHONY: upgrade
## Show outdated direct dependencies
upgrade:
	@echo "〉go mod upgrade"
	@rm -f go.work go.work.sum
	@$(foreach mod,$(GOMODS), (cd $(dir $(mod)) && echo "📂 $(dir $(mod))" && go list -u -m -f '{{if and (not .Indirect) .Update}}{{.Path}}{{end}}' all | xargs -n1 -I{} go get {}@latest) &&) true
	@$(MAKE) tidy

### Documentation

.PHONY: docs
## Open docs
docs:
	@echo "〉starting docs"
	@cd docs && bun install && bun run dev

.PHONY: docs.build
## Open docs
docs.build:
	@echo "〉building docs"
	@cd docs && bun install && bun run build

.PHONY: godocs
## Open go docs
godocs:
	@echo "〉starting go docs"
	@go doc -http

### Utils

.PHONY: help
# https://patorjk.com/software/taag/#p=display&f=Tmplr&t=testcontainers&x=none&v=4&h=4&w=80&we=false
## Show help text
help: g=\033[0;32m
help: b=\033[0;34m
help: w=\033[0;90m
help: e=\033[0m
help:
	@echo "$(g)"
	@echo "             •"
	@echo "╋┏┓┏╋┏┏┓┏┓╋┏┓┓┏┓┏┓┏┓┏"
	@echo "┗┗ ┛┗┗┗┛┛┗┗┗┻┗┛┗┗ ┛ ┛"
	@echo "with ❤ foomo by bestbytes"
	@echo "$(e)"
	@echo "$(b)Usage:$(e)\n  make [task]"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "  %-21s $(w)%s$(e)\n\n", cmd, help; help=""; \
			printf "$(b)\n%s:$(e)\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  %-21s $(w)%s$(e)\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        $(w)" help "$(e)\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)
	@echo ""

