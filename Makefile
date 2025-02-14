check_install: 
	which swagger || go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: check_install
	PATH=$(shell go env GOPATH)/bin:$(PATH) swagger generate spec -o ./swagger.yaml --scan-models