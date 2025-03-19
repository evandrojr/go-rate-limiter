package limiterstrategy

import (
	"strconv"
	"testing"

	"github.com/evandrojr/go-rate-limiter/configs"
	"github.com/stretchr/testify/assert"
)

func TestValidaPorToken(t *testing.T) {
	segundoRegistrado = int64(1)
	segundo := int64(1)
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"A": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	ip := "192.168.1.1"
	token := "A"
	v := validaAcesso(segundo, ip, token)
	assert.Nil(t, v)
	v = validaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), LIMITED_MESSAGE, "Erro segundo acesso")
	v = validaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), LIMITED_MESSAGE)
}

func TestValidaPorIP(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"A": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	ip := "192.168.1.1"
	token := "Outro"
	v := validaAcesso(segundo, ip, token)
	assert.Nil(t, v)
	v = validaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), LIMITED_MESSAGE)
	v = validaAcesso(segundo, ip, token)
	assert.Contains(t, v.Error(), LIMITED_MESSAGE)
}

func TestMesmoSegundoAcessoIpLiberado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"A": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundoRegistrado = int64(1)
	segundo := int64(1)
	configs.Config.IpMaxReqPerSecond = 1
	ip := "192.168.1.1"
	err := registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoIpBloqueado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"A": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	ip := "192.168.1.1"
	err := registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = registraAcessoIp(segundo, ip)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoTokenLiberado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	token := "token1"
	err := registraAcessoToken(segundo, token)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMesmoSegundoAcessoTokenBloqueado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	token := "token1"
	err := registraAcessoToken(segundo, token)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = registraAcessoToken(segundo, token)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestOutroSegundoAcessoIpLiberadoDentroCota(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(2)
	ip := "192.168.1.1"
	err := registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestOutroSegundoAcessoIpLiberado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	ip := "192.168.1.1"
	err := registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	segundo = int64(2)
	err = registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}

}

func TestOutroSegundoAcessoIpBloqueado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	ip := "192.168.1.1"
	err := registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	segundo = int64(2)
	err = registraAcessoIp(segundo, ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	err = registraAcessoIp(segundo, ip)
	if err == nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
}

func TestMultiplosAcessos(t *testing.T) {

	token1 := "token1"
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{token1: 3},
		IpMaxReqPerSecond:     3,
		BlockIpTime:           100,
		BlockTokenTime:        100,
	}
	Initialize(1, cfg)

	for i := 0; i < 200; i++ {
		go func(i int, t *testing.T) {
			err := valAcesso(i, t)
			if err != nil {
				t.Errorf("Erro ao registrar acesso: %s", err)
			}
		}(i, t)
		go func(i int, t *testing.T) {
			err := valAcesso(i, t)
			if err != nil {
				t.Errorf("Erro ao registrar acesso: %s", err)
			}
		}(i, t)
		go func(i int, t *testing.T) {
			err := valAcesso(i, t)
			if err != nil {
				t.Errorf("Erro ao registrar acesso: %s", err)
			}
		}(i, t)

	}
}

func valAcesso(i int, t *testing.T) error {
	ip := strconv.Itoa(i)
	err := validaAcesso(int64(1), ip, "token"+ip)
	if err != nil {
		t.Errorf("Erro ao registrar acesso: %s", err)
	}
	return nil
}
