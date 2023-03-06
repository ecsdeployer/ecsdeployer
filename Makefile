# gopkgs := ./cmd/... ./internal/... ./pkg/...
gopkgs := $(shell go list ./cmd/... ./internal/... ./pkg/... | grep -v testutil)


.PHONY: precommit
precommit: tidy lint test schema docs-pre

.PHONY: generate
generate:
	rm -f ./internal/fargate/sizes_gen.go
	go generate ./...

.PHONY: tidy
tidy:
	go mod verify
	go mod tidy
	@if ! git diff --quiet go.mod go.sum; then \
		echo "please run go mod tidy and check in changes, you might have to use the same version of Go as the CI"; \
		exit 1; \
	fi

.PHONY: lint-install
lint-install:
	@echo "Installing golangci-lint"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

.PHONY: lint
lint:
	@which golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found, please run: make lint-install"; \
		exit 1; \
	}
	@golangci-lint version
	@golangci-lint run && echo "Code passed lint check!"

.PHONY: test-release
test-release:
	goreleaser release --skip-publish --rm-dist --snapshot

.PHONY: check
check:
	go run . check --debug -c cmd/testdata/valid.yml
	go run . check --debug -c cmd/testdata/smoke.yml
	go run . check --debug -c www/docs/static/examples/generic.yml
	go run . check --debug -c www/docs/static/examples/simple_web.yml

.PHONY: schema
schema:
	@go run . schema -o ./www/docs/static/schema.json
	@cat ./www/docs/static/schema.json | jq -r .

.PHONY: smokedeploy-debug
smokedeploy-debug:
	@env AWS_PROFILE=ecsdeployer-example go run . deploy --debug -c cmd/testdata/smoke.yml --tag test --app-version 1.2.3

.PHONY: smokedeploy
smokedeploy:
	@env AWS_PROFILE=ecsdeployer-example go run . deploy -c cmd/testdata/smoke.yml --tag test --app-version 1.2.3

.PHONY: gen-man
gen-man:
	@./scripts/manpages.sh

.PHONY: showman
showman: gen-man
	@gunzip -c manpages/ecsdeployer.1.gz | nroff -man - | more -s

.PHONY: test
test:
	@./scripts/run_with_test_env.sh go test -timeout 180s $(gopkgs)

.PHONY: test-v
test-v:
	@./scripts/run_with_test_env.sh go test -v -timeout 180s $(gopkgs)

.PHONY: test-testutil
test-testutil:
	@go test -timeout 180s ./internal/testutil

.PHONY: test-single
test-single:
ifndef name
	$(error Rerun as 'make <command> name=<something>')
endif
	@./scripts/run_with_test_env.sh go test -v -timeout 180s $(gopkgs) -run $(name)

.PHONY: docs-serve
docs-serve:
	cd www && mkdocs serve

.PHONY: docs-deploy
docs-deploy: generate docs-pre
	cd www && mkdocs gh-deploy -c -b gh-pages -r newrepo --no-history

.PHONY: docs-pre
docs-pre:
	@./scripts/cmd_docs.sh

.PHONY: outdated
outdated:
	@go list -u -m -f '{{if not .Indirect}}{{if .Update}}{{.}}{{end}}{{end}}' all

.PHONY: coverage
coverage:
	@mkdir -p coverage
	
	@# This is only within the package itself
	@#./scripts/run_with_test_env.sh go test $(gopkgs) -cover -coverprofile=coverage/c.out -covermode=count

	@# will do coverage over the whole project
	@./scripts/run_with_test_env.sh go test $(gopkgs) -coverpkg=./... -coverprofile=coverage/c.out -covermode=count

	@# ignore testutil
	@cat coverage/c.out | grep -v ecsdeployer/internal/testutil/ | grep -v ecsdeployer/internal/builders/buildtestutils/ > coverage/c_notest.out
	@#go tool cover -html=coverage/c.out -o coverage/index.html
	@go tool cover -html=coverage/c_notest.out -o coverage/index.html

.PHONY: htmltest
htmltest:
	cd www && mkdocs build && htmltest -c htmltest.yml site