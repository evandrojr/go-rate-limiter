package limiter

import (
	"errors"
	"log"
	"time"

	"github.com/evandrojr/ratelimiter/configs"
)

var segundoRegistrado int64
var Acessos acessos
var tokenNotFound = "TOKEN_NOT_FOUND"

type acessos struct {
	Ip     map[string]int
	Tokens map[string]int
}

func Init() {
	Acessos = acessos{
		Ip:     make(map[string]int),
		Tokens: make(map[string]int),
	}
	segundoRegistrado = time.Now().Unix()
}

func ValidaAcesso(segundo int64, ip string, token string) error {
	RegistraAcessoTokenErr := RegistraAcessoToken(segundo, token)
	if RegistraAcessoTokenErr != nil && RegistraAcessoTokenErr.Error() == tokenNotFound {
		errRegistraAcessoIp := RegistraAcessoIp(segundo, ip)
		if errRegistraAcessoIp != nil {
			return errRegistraAcessoIp
		}
	}
	return RegistraAcessoTokenErr
}

func RegistraAcessoIp(segundo int64, ip string) error {
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Ip = make(map[string]int)
	}
	if Acessos.Ip[ip] > configs.Config.Ip {
		return errors.New("limite de acessos por IP excedido para o IP: " + ip)
	}
	Acessos.Ip[ip]++
	return nil
}

func RegistraAcessoToken(segundo int64, token string) error {
	if configs.Config.Tokens[token] == 0 {
		log.Println(tokenNotFound)
		return errors.New(tokenNotFound)
	}
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Tokens = make(map[string]int)
	}
	if Acessos.Tokens[token] > configs.Config.Tokens[token] {
		return errors.New("limite de acessos por token excedido para o token: " + token)
	}
	Acessos.Tokens[token]++
	return nil
}
