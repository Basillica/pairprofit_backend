.PHONY: run-debug

GIT_COMMIT := $(shell git rev-list -1 HEAD)
run-debug:
	export GIN_MODE=debug && go build -ldflags="-X 'main.Version=$(GIT_COMMIT)'" -o pairprofit_backend && ./pairprofit_backend

run-release:
	export GIN_MODE=release && go build -ldflags="-X 'main.Version=$(GIT_COMMIT)'" -o pairprofit_backend && ./pairprofit_backend

lint:
	golangci-lint run
