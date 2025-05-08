#!/usr/bin/env bash

echo "Banco de dados pronto."
# cd cmd/ordersystem && go run main.go wire_gen.go
echo "Subindo a aplicação"
DOCKER_EXECUTION=true go run app.go 
# echo "Executando testes unitários com verificação de concorrência..."
# go test -race ./... 
# echo "Executando testes de integração..."
# ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" http://localhost:8080/
# ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" -H "APT_KEY: DESCONHECIDO" http://localhost:8080/
# ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" -H "APT_KEY: TRAVAEU" http://localhost:8080/
