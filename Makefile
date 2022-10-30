build:
	go build -o bin/currency-conversion  ./cmd/currency-conversion

run:
	./bin/currency-conversion

start:
	clear && go build -o bin/currency-conversion  ./cmd/currency-conversion && ./bin/currency-conversion

clean:
	rm -rf bin && rm -rf internal/mocks && clear

migrate-tables:
	go build -o postgres_migrate  ./cmd/seed  && ./postgres_migrate -type=migrate && rm -rf postgres_migrate && clear

drop-tables:
	go build -o postgres_drop  ./cmd/seed  && ./postgres_drop -type=drop && rm -rf postgres_drop && clear

generate-mocks:
	mockgen -destination=internal/mocks/auth/mock_user_repository.go -package mocks github.com/eneskzlcn/currency-conversion-service/internal/auth AuthRepository
	mockgen -destination=internal/mocks/auth/mock_auth_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/internal/auth AuthService
	mockgen -destination=internal/mocks/exchange/mock_exchange_service.go -package mocks github.com/eneskzlcn/currency-conversion-service/internal/exchange ExchangeService
