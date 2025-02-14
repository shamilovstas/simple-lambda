.PHONY: clean bootstrap deploy
bootstrap:
	@GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o $@ cmd/main.go

deploy: bootstrap
	@terraform apply -auto-approve

clean:
	@rm -f bootstrap payload.zip
