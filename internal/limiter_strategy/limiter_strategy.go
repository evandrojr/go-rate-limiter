package limiterstrategy

import (
	"github.com/evandrojr/go-rate-limiter/configs"
	persistencystrategy "github.com/evandrojr/go-rate-limiter/internal/persistency_strategy"
)

type LimiterStrategyStruct struct{}

type LimiterStrategyI interface {
	Init(segundoRegistrado int64, configs configs.EnvConfig, logger persistencystrategy.PersistencyStrategyI)
	ValidaAcesso(segundoRegistrado int64, ip string, token string) error
}

func ValidaAcessoPolimorfico(estrategia LimiterStrategyI, segundoRegistrado int64, ip string, token string) error {
	return estrategia.ValidaAcesso(segundoRegistrado, ip, token)
}

type TipoEstrategiaStruct struct {
	estrategy LimiterStrategyI
}

// Permite definir a estratégia a ser utilizada.
func (e *TipoEstrategiaStruct) SetStrategy(strategy LimiterStrategyI) {
	e.estrategy = strategy
}

// Permite ler a estratégia utilizada.
func (e *TipoEstrategiaStruct) GetStrategy() LimiterStrategyI {
	return e.estrategy
}
