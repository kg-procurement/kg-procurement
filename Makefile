build-ci:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	@echo ">> Building Main Core Binary..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o ./bin/kg-procurement ./cmd
	@echo ">> Finished"

tidy:
	@go mod tidy

tool-mock:
	mkdir -p ./bin
	GOBIN=$(CURDIR)/bin go install go.uber.org/mock/mockgen@latest && go get go.uber.org/mock/gomock

gen-mock: tool-mock
	PATH=$(CURDIR)/bin:${PATH} go generate ./...

docker-up:
	docker-compose -f docker-compose.yaml up -d
docker-down:
	docker-compose -f docker-compose.yaml down
