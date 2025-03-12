package limiterstrategy

type LimiterStrategyStruct struct{}

type LimiterStrategyI interface {
	ValidaAcesso(segundoRegistrado int64, ip string, token string) error
}

func ValidaAcessoPolimorfico(estrategia LimiterStrategyI, segundoRegistrado int64, ip string, token string) error {
	return estrategia.ValidaAcesso(segundoRegistrado, ip, token)
}

//)

// func (l limiterstrategy) ValidaAcesso(segundoRegistrado int64, ip string, token string) error {

// }

// func NewLimiterStrategy() LimiterStrategy {
// 	return &limiterstrategy{}
// }
