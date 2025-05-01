#!/usr/bin/env bash

# echo "Aguardando banco e migrações..."
# until ./migrate -path=./internal/infra/database/migrations -database "mysql://root:root@tcp(mysqldb:3306)/orders" -verbose up; do
#   sleep 2
# done

echo "Banco de dados pronto. Iniciando aplicação..."
# cd cmd/ordersystem && go run main.go wire_gen.go
echo "Rodando testes unitários com mais de 300k requisições e iniciando aplicação em caso de sucesso..."
DOCKER_EXECUTION=true go test -race ./... && DOCKER_EXECUTION=true go run app.go
# DOCKER_EXECUTION=true go run app.go