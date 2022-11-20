.PHONY: build-docker
build-docker:
	docker build . --tag filesystem-api

.PHONY: run-docker
run-docker:
	docker run -p 8080:8080 filesystem-api

.PHONY: test
test:
	go test ./...
