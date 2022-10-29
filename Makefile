build:
	go build -o bin/currency-conversion  ./cmd/currency-conversion

run:
	./bin/currency-conversion

start:
	clear && go build -o bin/currency-conversion  ./cmd/currency-conversion && ./bin/currency-conversion

clean:
	rm -rf bin && clear

migrate-tables:
	go build -o postgres_migrate  ./cmd/seed  && ./postgres_migrate -type=migrate && rm -rf postgres_migrate && clear

drop-tables:
	go build -o postgres_drop  ./cmd/seed  && ./postgres_drop -type=drop && rm -rf postgres_drop && clear