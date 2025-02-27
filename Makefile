.PHONY: all
all: test check coverage build

.PHONY: build
build: mp4ff-crop mp4ff-decrypt mp4ff-encrypt mp4ff-info mp4ff-nallister mp4ff-pslister mp4ff-subslister examples

.PHONY: prepare
prepare:
	go mod vendor

mp4ff-crop mp4ff-decrypt mp4ff-encrypt mp4ff-info mp4ff-nallister mp4ff-pslister mp4ff-subslister:
	go build -ldflags "-X github.com/vtpl1/mp4ff/mp4.commitVersion=$$(git describe --tags --always) -X github.com/vtpl1/mp4ff/mp4.commitDate=$$(git log -1 --format=%ct)" -o out/$@ ./cmd/$@/main.go

.PHONY: examples
examples: add-sidx combine-segs initcreator multitrack resegmenter segmenter

add-sidx combine-segs initcreator multitrack resegmenter segmenter:
	go build -o examples-out/$@  ./examples/$@

.PHONY: test
test: prepare
	go test -cover ./...

.PHONY: testsum
testsum: prepare
	gotestsum

.PHONY: open-docs
open-docs:
	echo "If needed: go install golang.org/x/pkgsite/cmd/pkgsite@latest"
	pkgsite -http localhost:9999
	# open http://localhost:9999/pkg/github.com/vtpl1/mp4ff/

.PHONY: coverage
coverage:
	# Ignore (allow) packages without any tests
	set -o pipefail
	go test ./... -coverprofile coverage.out
	set +o pipefail
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out -o coverage.txt
	tail -1 coverage.txt

.PHONY: check
check: prepare
	golangci-lint run

clean:
	rm -f out/*
	rm -r examples-out/*

install:
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install mvdan.cc/gofumpt@latest
	@go install golang.org/x/lint/golint@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install gotest.tools/gotestsum@latest

lintverify:
	@golangci-lint config verify

lint:
	@golangci-lint run ./...
	
install1: all
	cp out/* $(GOPATH)/bin/

.PHONY: fmt
fmt:
	@gofumpt -l -w .
