TEST?=./infoblox
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

default: test

bin: generate
	@sh -c "'$(CURDIR)/scripts/build.sh'"

ci: generate
		@sh -c "'$(CURDIR)/scripts/build-ci.sh'"

# get dependencies
updatedeps:
	go list ./... \
        | xargs go list -f '{{join .Deps "\n"}}' \
		| grep -v terraform-provider-infoblox \
        | grep -v '/internal/' \
        | sort -u \
        | xargs go get -f -u -v

# test runs the unit tests and vets the code
test: generate
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	@$(MAKE) vet

# testacc runs acceptance tests
testacc: generate
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 90m

# testrace runs the race checker
testrace: generate
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS) $(TEST) "
	@go tool vet $(VETARGS) $(TEST) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	go generate $$(go list ./... | grep -v /vendor/)

.PHONY: bin default generate test updatedeps vet
