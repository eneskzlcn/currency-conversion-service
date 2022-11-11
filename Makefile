build:
	go build -o bin/currency-conversion

run:
	./bin/currency-conversion

start:
	swag init && go build -o bin/currency-conversion  && ./bin/currency-conversion

clean:
	rm -rf bin && rm -rf app/mocks && rm -rf postgres_migrate && rm -rf postgres_drop && rm -rf unit_coverage.out && clear

unit-tests:
	go test -v ./... -coverprofile=unit_coverage.out -short -tags=unit

linter:
	golangci-lint run
migrate-tables:
	go build -o postgres_migrate  ./seed/cmd  && ./postgres_migrate -type=migrate && rm -rf postgres_migrate && clear

drop-tables:
	go build -o postgres_drop  ./seed/cmd  && ./postgres_drop -type=drop && rm -rf postgres_drop && clear

swagger:
	swag init

generate-mocks:
	mockgen -destination=app/mocks/auth/mock_auth_repository.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/auth Repository
	mockgen -destination=app/mocks/auth/mock_auth_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/auth Service
	mockgen -destination=app/mocks/exchange/mock_exchange_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/exchange Service
	mockgen -destination=app/mocks/exchange/mock_auth_guard.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/exchange AuthGuard
	mockgen -destination=app/mocks/exchange/mock_exchange_repository.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/exchange Repository
	mockgen -destination=app/mocks/conversion/mock_conversion_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/conversion Service
	mockgen -destination=app/mocks/conversion/mock_auth_guard.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/conversion AuthGuard
	mockgen -destination=app/mocks/conversion/mock_wallet_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/conversion WalletService
	mockgen -destination=app/mocks/conversion/mock_rabbitmq_producer.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/conversion RabbitmqProducer
	mockgen -destination=app/mocks/conversion/mock_conversion_repository.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/conversion Repository
	mockgen -destination=app/mocks/wallet/mock_wallet_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/wallet Service
	mockgen -destination=app/mocks/wallet/mock_auth_guard.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/wallet AuthGuard
	mockgen -destination=app/mocks/wallet/mock_wallet_repository.go -package mocks github.com/eneskzlcn/currency-conversion-service/app/wallet Repository
