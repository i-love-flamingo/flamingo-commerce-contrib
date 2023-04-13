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
	go test -test.count=10 -race -v ./cart/redis/integrationtest/... -tags=integration

generate-integrationtest-graphql:
	rm -f cart/redis/integrationtest/graphql/generated.go
	rm -f cart/redis/integrationtest/graphql/resolver.go
	go generate ./...
	export RUN="0" && cd cart/redis/integrationtest && go run -tags graphql main.go

fix:
	gofmt -l -w .
