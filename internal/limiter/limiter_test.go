package limiter

import (
	"testing"

	"github.com/evandrojr/ratelimiter/configs"
)

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
