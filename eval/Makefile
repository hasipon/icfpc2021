
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build
	GOOS=darwin GOARCH=amd64 go build -o eval.darwin
	GOOS=windows GOARCH=amd64 go build
.PHONY: test
test:
	bash test/test.sh || echo failed