FROM golang:1.24.0

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y apache2-utils

RUN go mod tidy

CMD ["sh", "-c", "go test -race ./... "]

