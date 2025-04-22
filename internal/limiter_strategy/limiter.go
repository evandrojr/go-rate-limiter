package limiterstrategy

import (
	"errors"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/evandrojr/go-rate-limiter/configs"
	"github.com/evandrojr/go-rate-limiter/internal/redis"
	"github.com/kr/pretty"

	"sync"
)

var acessosMutex sync.Mutex
var segundoRegistradoMutex sync.Mutex

var _segundoRegistrado int64
var _acessos acessosType
var tokenNotFound = "TOKEN_NOT_FOUND"
var exceedIpLimit = "limite de acessos por IP excedido para o IP: "
var exceedTokenLimit = "limite de acessos por token excedido para o token: "

const LIMITED_MESSAGE = "you have reached the maximum number of requests or actions allowed within a certain time frame"

type acessoRecord struct {
	NumeroAcessosNoSegundo int
	BloqueadoAte           int64
}

type acessosType struct {
	Ip     map[string]acessoRecord
	Tokens map[string]acessoRecord
}

var envConfig configs.EnvConfig

func (l LimiterStrategyStruct) Init(segundoRegistrado int64, configs configs.EnvConfig) {
	Initialize(segundoRegistrado, configs)
}

func Initialize(secReg int64, configs configs.EnvConfig) {
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	envConfig = configs
	_acessos = (acessosType{
		Ip:     make(map[string]acessoRecord),
		Tokens: make(map[string]acessoRecord),
	})
	atomic.StoreInt64(&_segundoRegistrado, secReg)
}

func bloquearIp(segundoAtual int64, ip string) {
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	_acessos.Ip[ip] = acessoRecord{
		NumeroAcessosNoSegundo: _acessos.Ip[ip].NumeroAcessosNoSegundo,
		BloqueadoAte:           segundoAtual + int64(envConfig.BlockIpTime),
	}
}

func bloquearToken(segundoAtual int64, token string) {
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	_acessos.Tokens[token] = acessoRecord{
		NumeroAcessosNoSegundo: _acessos.Tokens[token].NumeroAcessosNoSegundo,
		BloqueadoAte:           segundoAtual + int64(envConfig.BlockTokenTime),
	}
}

func verificaBloqueio(segundoAtual int64, ip string, token string) error {
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	if _acessos.Ip[ip].BloqueadoAte >= segundoAtual {
		return errors.New("IP bloqueado até " + strconv.FormatInt(_acessos.Ip[ip].BloqueadoAte, 10))
	}
	if _acessos.Tokens[token].BloqueadoAte >= segundoAtual {
		return errors.New("Token bloqueado até " + strconv.FormatInt(_acessos.Tokens[token].BloqueadoAte, 10))
	}
	return nil
}

func (l LimiterStrategyStruct) ValidaAcesso(segundoRegistrado int64, ip string, token string) error {
	error := validaAcesso(segundoRegistrado, ip, token)
	if error != nil {
		redis.RPush("segundoRegistrado: " + strconv.FormatInt(segundoRegistrado, 10) + " ip: " + ip + " token: " + token)
		log.Println(error)
	}
	return error
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
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	segundoRegistradoMutex.Lock()
	defer segundoRegistradoMutex.Unlock()
	if segundo != _segundoRegistrado {
		_segundoRegistrado = segundo
		_acessos.Ip[ip] = acessoRecord{}
	}
	pretty.Println(_acessos.Ip[ip])
	if _acessos.Ip[ip].NumeroAcessosNoSegundo >= envConfig.IpMaxReqPerSecond {
		return errors.New(exceedIpLimit + ip)
	}
	ar := acessoRecord{
		NumeroAcessosNoSegundo: _acessos.Ip[ip].NumeroAcessosNoSegundo + 1,
		BloqueadoAte:           _acessos.Ip[ip].BloqueadoAte}
	_acessos.Ip[ip] = ar
	return nil
}

func registraAcessoToken(segundo int64, token string) error {
	acessosMutex.Lock()
	defer acessosMutex.Unlock()
	segundoRegistradoMutex.Lock()
	defer segundoRegistradoMutex.Unlock()
	if envConfig.TokensMaxReqPerSecond[token] == 0 {
		log.Println(tokenNotFound)
		return errors.New(tokenNotFound)
	}
	pretty.Println(_acessos.Tokens[token])
	if segundo != _segundoRegistrado {
		_segundoRegistrado = segundo
		_acessos.Tokens[token] = acessoRecord{}
	}
	if _acessos.Tokens[token].NumeroAcessosNoSegundo >= envConfig.TokensMaxReqPerSecond[token] {
		return errors.New(exceedTokenLimit + token)
	}
	ar := acessoRecord{
		NumeroAcessosNoSegundo: _acessos.Tokens[token].NumeroAcessosNoSegundo + 1,
		BloqueadoAte:           _acessos.Tokens[token].BloqueadoAte}
	_acessos.Tokens[token] = ar
	return nil
}
