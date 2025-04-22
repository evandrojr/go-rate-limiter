# Rate Limiter Go

Os requisitos estão em [Requisitos.md](Requisitos.md)

##

Execução:

`docker-compose up`

## 

Testes unitários:

go test -race -v  ./...

##

Ver os log no Redis

`
./redis-console.sh
LRANGE limits:log 0 -1
`