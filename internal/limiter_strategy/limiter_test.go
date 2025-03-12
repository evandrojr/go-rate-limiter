package limiterstrategy

import (
	"testing"

	"github.com/evandrojr/go-rate-limiter/configs"
	"github.com/stretchr/testify/assert"
)

func TestValidaPorToken(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	configs.Config.Tokens = map[string]int{"A": 1}
	ip := "192.168.1.1"
	token := "A"
	v := ValidaAcesso(segundo, ip, token)
	assert.Nil(t, v)
	v = ValidaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), exceedTokenLimit, "Erro segundo acesso")
	v = ValidaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), exceedTokenLimit)
}

func TestValidaPorIP(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	configs.Config.Tokens = map[string]int{"A": 10}
	ip := "192.168.1.1"
	token := "Outro"
	v := ValidaAcesso(segundo, ip, token)
	assert.Nil(t, v)
	v = ValidaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), exceedIpLimit)
	v = ValidaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), exceedIpLimit)
}

func TestMesmoSegundoAcessoIpLiberado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	ip := "192.168.1.1"
	err := RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoIpBloqueado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	ip := "192.168.1.1"
	err := RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = RegistraAcessoIp(segundo, ip)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoTokenLiberado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Tokens = map[string]int{"token1": 1}
	token := "token1"
	err := RegistraAcessoToken(segundo, token)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoTokenBloqueado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Tokens = map[string]int{"token1": 1}
	token := "token1"
	err := RegistraAcessoToken(segundo, token)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = RegistraAcessoToken(segundo, token)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestOutroSegundoAcessoIpLiberadoDentroCota(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(2)
	configs.Config.Ip = 1
	ip := "192.168.1.1"
	err := RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestOutroSegundoAcessoIpLiberado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	ip := "192.168.1.1"
	err := RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	segundo = int64(2)
	err = RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}

}

func TestOutroSegundoAcessoIpBloqueado(t *testing.T) {
	Init()
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.Ip = 1
	ip := "192.168.1.1"
	err := RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	segundo = int64(2)
	err = RegistraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = RegistraAcessoIp(segundo, ip)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

// func TestMesmoSegundoAcessoTokenLiberado(t *testing.T) {
// 	Init()
// 	segundoRegistrado = int64(1)
// 	segundo := int64(1)
// 	configs.Config.Tokens = map[string]int{"token1": 1}
// 	token := "token1"
// 	err := RegistraAcessoToken(segundo, token)
// 	if err != nil {
// 		t.Errorf("Erro ao registrar acesso: %s", err)
// 	}
// }

// func TestMesmoSegundoAcessoTokenBloqueado(t *testing.T) {
// 	Init()
// 	segundoRegistrado = int64(1)
// 	segundo := int64(1)
// 	configs.Config.Tokens = map[string]int{"token1": 1}
// 	token := "token1"
// 	err := RegistraAcessoToken(segundo, token)
// 	if err != nil {
// 		t.Errorf("Erro ao registrar acesso: %s", err)
// 	}
// 	err = RegistraAcessoToken(segundo, token)
// 	if err == nil {
// 		t.Errorf("Erro ao registrar acesso: %s", err)
// 	}
// }
