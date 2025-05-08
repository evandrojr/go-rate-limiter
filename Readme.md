# Rate Limiter Go

Os requisitos estão em [Requisitos.md](Requisitos.md)

## Execução:

`docker-compose up`

## Testes unitários:

São executados antes do inicío da aplicação mais 300k requisições 

go test -race  ./...

## Testes com apache AB

Rode com docker-compose up -d senão a escrita de log no terminal vai deixar os testes muito lentos
```


sudo apt-get update && sudo apt-get install -y apache2-utils

ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" http://localhost:8080/

ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" -H "APT_KEY: DESCONHECIDO" http://localhost:8080/

ab -n 1000 -c 10 -H "Host: localhost:8080" -H "Content-Type: application/json" -H "APT_KEY: TRAVAEU" http://localhost:8080/
```

Ver os log no Redis

`
./redis-console.sh
LRANGE limits:log 0 -1
`