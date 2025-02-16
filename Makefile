.PHONY: clean bootstrap deploy test
bootstrap:
	@GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o $@ cmd/main.go

deploy: bootstrap
	@terraform apply -auto-approve

test:
	@go test -v ./...

clean:
	@rm -f bootstrap payload.zip
