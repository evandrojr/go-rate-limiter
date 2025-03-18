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

var envConfig configs.EnvConfig

func Init(configs configs.EnvConfig) {
	envConfig = configs
	Acessos = acessos{
		Ip:     make(map[string]acessoRecord),
		Tokens: make(map[string]acessoRecord),
	}
	segundoRegistrado = time.Now().Unix()
}

func BloquearIp(ip string) {
	Acessos.Ip[ip] = acessoRecord{
		NumeroAcessosNoSegundo: 0,
		BloqueadoAte:           time.Now().Add(time.Second * time.Duration(envConfig.BlockIpTime)),
	}
}

func BloquearToken(token string) {
	Acessos.Tokens[token] = acessoRecord{
		NumeroAcessosNoSegundo: 0,
		BloqueadoAte:           time.Now().Add(time.Second * time.Duration(envConfig.BlockTokenTime)),
	}
}

func VerificaBloqueio(segundo int64, ip string, token string) error {
	if Acessos.Ip[ip].BloqueadoAte.After(time.Now()) {
		return errors.New("IP bloqueado até " + Acessos.Ip[ip].BloqueadoAte.String())
	}
	if Acessos.Tokens[token].BloqueadoAte.After(time.Now()) {
		return errors.New("Token bloqueado até " + Acessos.Tokens[token].BloqueadoAte.String())
	}
	return nil
}

func (l LimiterStrategyStruct) ValidaAcesso(segundoRegistrado int64, ip string, token string) error {
	return ValidaAcesso(segundoRegistrado, ip, token)
}

func ValidaAcesso(segundo int64, ip string, token string) error {
	if err := VerificaBloqueio(segundo, ip, token); err != nil {
	}

	RegistraAcessoTokenErr := RegistraAcessoToken(segundo, token)
	if RegistraAcessoTokenErr != nil && RegistraAcessoTokenErr.Error() == tokenNotFound {
		return RegistraAcessoIp(segundo, ip)
	}
	return RegistraAcessoTokenErr
}

func RegistraAcessoIp(segundo int64, ip string) error {
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Ip = make(map[string]acessoRecord)
	}
	println("Acessos.Ip[ip]", Acessos.Ip[ip])
	if Acessos.Ip[ip].NumeroAcessosNoSegundo >= envConfig.IpMaxReqPerSecond {
		return errors.New(exceedIpLimit + ip)
	}
	Acessos.Ip[ip].++
	return nil
}

func RegistraAcessoToken(segundo int64, token string) error {
	if envConfig.TokensMaxReqPerSecond[token] == 0 {
		log.Println(tokenNotFound)
		return errors.New(tokenNotFound)
	}
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Tokens = make(map[string]int)
	}
	if Acessos.Tokens[token] >= envConfig.TokensMaxReqPerSecond[token] {
		return errors.New(exceedTokenLimit + token)
	}
	Acessos.Tokens[token]++
	return nil
}
