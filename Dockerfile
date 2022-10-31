FROM    golang:1.18-alpine3.16
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod tidy -go=1.18
RUN go build -o bin/currency-conversion ./cmd/currency-conversion
EXPOSE 4001
CMD [ "bin/currency-conversion" ]