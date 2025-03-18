package limiterstrategy

import (
	"errors"
	"log"
	"time"

	"github.com/evandrojr/go-rate-limiter/configs"
)

var segundoRegistrado int64
var Acessos acessos
var tokenNotFound = "TOKEN_NOT_FOUND"
var exceedIpLimit = "limite de acessos por IP excedido para o IP: "
var exceedTokenLimit = "limite de acessos por token excedido para o token: "

type acessoRecord struct {
	NumeroAcessosNoSegundo int
	BloqueadoAte           time.Time
}

type acessos struct {
	Ip     map[string]acessoRecord
	Tokens map[string]acessoRecord
}

func Init(configs configs.Config) {
	Acessos = acessos{
		Ip:     make(map[string]acessoRecord),
		Tokens: make(map[string]acessoRecord),
	}
	segundoRegistrado = time.Now().Unix()
}

func BloquearIp(ip string) {

}

func (l LimiterStrategyStruct) ValidaAcesso(segundoRegistrado int64, ip string, token string) error {
	return ValidaAcesso(segundoRegistrado, ip, token)
}

func ValidaAcesso(segundo int64, ip string, token string) error {
	RegistraAcessoTokenErr := RegistraAcessoToken(segundo, token)
	if RegistraAcessoTokenErr != nil && RegistraAcessoTokenErr.Error() == tokenNotFound {
		return RegistraAcessoIp(segundo, ip)
	}
	return RegistraAcessoTokenErr
}

func RegistraAcessoIp(segundo int64, ip string) error {
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Ip = make(map[string]int)
	}
	println("Acessos.Ip[ip]", Acessos.Ip[ip])
	if Acessos.Ip[ip] >= configs.Config.IpMaxReqPerSecond {
		return errors.New(exceedIpLimit + ip)
	}
	Acessos.Ip[ip]++
	return nil
}

func RegistraAcessoToken(segundo int64, token string) error {
	if configs.Config.TokensMaxReqPerSecond[token] == 0 {
		log.Println(tokenNotFound)
		return errors.New(tokenNotFound)
	}
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Tokens = make(map[string]int)
	}
	if Acessos.Tokens[token] >= configs.Config.TokensMaxReqPerSecond[token] {
		return errors.New(exceedTokenLimit + token)
	}
	Acessos.Tokens[token]++
	return nil
}
