package limiterstrategy

import (
	"strconv"
	"sync"
	"testing"

	"github.com/evandrojr/go-rate-limiter/configs"
	"github.com/stretchr/testify/assert"
)

func TestValidaPorToken(t *testing.T) {
	_segundoRegistrado = int64(1)
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
	_segundoRegistrado = int64(1)
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
		IpMaxReqPerSecond:     3000,
		BlockIpTime:           1000,
		BlockTokenTime:        1000,
	}
	Initialize(1, cfg)
	qtd := 1000

	var wg sync.WaitGroup
	for i := 0; i <= qtd; i++ {
		wg.Add(3) // Adiciona 3 goroutines ao WaitGroup

		go func(i int, t *testing.T) {
			defer wg.Done() // Marca a goroutine como concluída
			valAcesso(i)
		}(i, t)

		go func(i int, t *testing.T) {
			defer wg.Done()
			valAcesso(i)
		}(i, t)

		go func(i int, t *testing.T) {
			defer wg.Done()
			valAcesso(i)
		}(i, t)
	}

	wg.Wait() // Aguarda todas as goroutines serem concluídas
}

func TestValidaAcessoTokenBloqueado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 1},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	token := "token1"
	ip := "192.168.1.1"

	// First access should pass
	err := validaAcesso(segundo, ip, token)
	assert.Nil(t, err)

	// Second access should block the token
	err = validaAcesso(segundo, ip, token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), LIMITED_MESSAGE)
}

func TestValidaAcessoIpBloqueado(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{"token1": 0},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	token := "token1"
	ip := "192.168.1.1"

	// First access should pass
	err := validaAcesso(segundo, ip, token)
	assert.Nil(t, err)

	// Second access should block the IP
	err = validaAcesso(segundo, ip, token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), LIMITED_MESSAGE)
}

func TestValidaAcessoTokenNotFound(t *testing.T) {
	cfg := configs.EnvConfig{
		TokensMaxReqPerSecond: map[string]int{},
		IpMaxReqPerSecond:     1,
		BlockIpTime:           2,
		BlockTokenTime:        2,
	}
	Initialize(1, cfg)
	segundo := int64(1)
	token := "unknown_token"
	ip := "192.168.1.1"

	// Access with an unknown token should fallback to IP validation
	err := validaAcesso(segundo, ip, token)
	assert.Nil(t, err)

	// Second access should block the IP
	err = validaAcesso(segundo, ip, token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), LIMITED_MESSAGE)
}

func valAcesso(i int) {
	ip := strconv.Itoa(i)
	validaAcesso(int64(1), ip, "token"+ip)

}
