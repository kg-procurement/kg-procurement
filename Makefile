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

seed-product:
	@go run scripts/seeder/main.go product
seed-product-category:
	@go run scripts/seeder/main.go product_category
seed-product-type:
	@go run scripts/seeder/main.go product_type
seed-uom:
	@go run scripts/seeder/main.go uom
seed-vendor:
	@go run scripts/seeder/main.go vendor
seed-product-vendor:
	@go run scripts/seeder/main.go product-vendor


migrate-up:
	@goose -dir migrations postgres "password=postgres user=postgres port=5432 dbname=kg-procurement host=localhost sslmode=disable" up
migrate-down:
	@goose -dir migrations postgres "password=postgres user=postgres port=5432 dbname=kg-procurement host=localhost sslmode=disable" down
migrate-status:
	@goose -dir migrations postgres "password=postgres user=postgres port=5432 dbname=kg-procurement host=localhost sslmode=disable" status