package limiterstrategy

import (
	"errors"
	"log"
	"strconv"

	"github.com/evandrojr/go-rate-limiter/configs"
	"github.com/kr/pretty"
)

var segundoRegistrado int64
var Acessos acessos
var tokenNotFound = "TOKEN_NOT_FOUND"
var exceedIpLimit = "limite de acessos por IP excedido para o IP: "
var exceedTokenLimit = "limite de acessos por token excedido para o token: "

const LIMITED_MESSAGE = "you have reached the maximum number of requests or actions allowed within a certain time frame"

type acessoRecord struct {
	NumeroAcessosNoSegundo int
	BloqueadoAte           int64
}

type acessos struct {
	Ip     map[string]acessoRecord
	Tokens map[string]acessoRecord
}

var envConfig configs.EnvConfig

func (l LimiterStrategyStruct) Init(segundoRegistrado int64, configs configs.EnvConfig) {
	Initialize(segundoRegistrado, configs)
}

func Initialize(secReg int64, configs configs.EnvConfig) {
	envConfig = configs
	Acessos = acessos{
		Ip:     make(map[string]acessoRecord),
		Tokens: make(map[string]acessoRecord),
	}
	segundoRegistrado = secReg
}

func bloquearIp(segundoAtual int64, ip string) {
	Acessos.Ip[ip] = acessoRecord{
		NumeroAcessosNoSegundo: Acessos.Ip[ip].NumeroAcessosNoSegundo,
		BloqueadoAte:           segundoAtual + int64(envConfig.BlockIpTime),
	}
}

func bloquearToken(segundoAtual int64, token string) {
	Acessos.Tokens[token] = acessoRecord{
		NumeroAcessosNoSegundo: Acessos.Tokens[token].NumeroAcessosNoSegundo,
		BloqueadoAte:           segundoAtual + int64(envConfig.BlockTokenTime),
	}
}

func verificaBloqueio(segundoAtual int64, ip string, token string) error {
	if Acessos.Ip[ip].BloqueadoAte >= segundoAtual {
		return errors.New("IP bloqueado até " + strconv.FormatInt(Acessos.Ip[ip].BloqueadoAte, 10))
	}
	if Acessos.Tokens[token].BloqueadoAte >= segundoAtual {
		return errors.New("Token bloqueado até " + strconv.FormatInt(Acessos.Tokens[token].BloqueadoAte, 10))
	}
	return nil
}

func (l LimiterStrategyStruct) ValidaAcesso(segundoRegistrado int64, ip string, token string) error {
	return validaAcesso(segundoRegistrado, ip, token)
}

func validaAcesso(segundo int64, ip string, token string) error {

	if err := verificaBloqueio(segundo, ip, token); err != nil {
		log.Println(err)
		return errors.New(LIMITED_MESSAGE)
	}
	RegistraAcessoTokenErr := registraAcessoToken(segundo, token)
	if RegistraAcessoTokenErr != nil && RegistraAcessoTokenErr.Error() == tokenNotFound {
		RegistraAcessoIpErro := registraAcessoIp(segundo, ip)
		if RegistraAcessoIpErro != nil {
			log.Println(RegistraAcessoIpErro)
			bloquearIp(segundo, ip)
			return errors.New(LIMITED_MESSAGE)
		}
		return nil
	}
	log.Println(RegistraAcessoTokenErr)
	bloquearToken(segundo, token)
	return RegistraAcessoTokenErr
}

func registraAcessoIp(segundo int64, ip string) error {
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Ip[ip] = acessoRecord{}
	}
	pretty.Println(Acessos.Ip[ip])
	if Acessos.Ip[ip].NumeroAcessosNoSegundo >= envConfig.IpMaxReqPerSecond {
		return errors.New(exceedIpLimit + ip)
	}
	ar := acessoRecord{
		NumeroAcessosNoSegundo: Acessos.Ip[ip].NumeroAcessosNoSegundo + 1,
		BloqueadoAte:           Acessos.Ip[ip].BloqueadoAte}
	Acessos.Ip[ip] = ar
	return nil
}

func registraAcessoToken(segundo int64, token string) error {
	if envConfig.TokensMaxReqPerSecond[token] == 0 {
		log.Println(tokenNotFound)
		return errors.New(tokenNotFound)
	}
	pretty.Println(Acessos.Tokens[token])
	if segundo != segundoRegistrado {
		segundoRegistrado = segundo
		Acessos.Tokens[token] = acessoRecord{}
	}
	if Acessos.Tokens[token].NumeroAcessosNoSegundo >= envConfig.TokensMaxReqPerSecond[token] {
		return errors.New(exceedTokenLimit + token)
	}
	ar := acessoRecord{
		NumeroAcessosNoSegundo: Acessos.Tokens[token].NumeroAcessosNoSegundo + 1,
		BloqueadoAte:           Acessos.Tokens[token].BloqueadoAte}
	Acessos.Tokens[token] = ar
	return nil
}
