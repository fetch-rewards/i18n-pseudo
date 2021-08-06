MAKE=make
GOTEST=go test
COVERPROFILE=coverage.out

export

.PHONY: test
test:
	$(GOTEST) -coverpkg=./... -coverprofile=$(COVERPROFILE) -outputdir=. -v -test.short ./...

linter:
	golangci-lint run

tidy:
	go mod tidy
