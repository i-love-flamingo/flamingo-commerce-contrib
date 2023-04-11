CONTEXT?=dev
REPLACE?=-replace flamingo.me/flamingo/v3=../flamingo -replace flamingo.me/form=../form
DROPREPLACE?=-dropreplace flamingo.me/flamingo/v3 -dropreplace flamingo.me/form

.PHONY: local unlocal test

test:
	go test -race -v ./...
	gofmt -l -e -d .
	golint ./...
	misspell -error .
	ineffassign .

integrationtest:
	go test -test.count=10 -race -v ./test/integrationtest/... -tags=integration

generate-integrationtest-graphql:
	rm -f test/integrationtest/projecttest/graphql/generated.go
	rm -f test/integrationtest/projecttest/graphql/resolver.go
	go generate ./...
	export RUN="0" && cd test/integrationtest/projecttest && go run -tags graphql main.go

fix:
	gofmt -l -w .

run-integrationtest-demo-project:
	cd test/integrationtest/projecttest/tests && RUN=1 INTEGRATION_TEST_PORT=10000 go run ../main.go
