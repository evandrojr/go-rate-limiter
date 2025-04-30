package persistencystrategy

type PersistencyStrategyStruct struct{}

type PersistencyStrategyI interface {
	Init(config []string)
	Log(msg string) error
}

func LogPolimorfico(estrategia PersistencyStrategyI, msg string) error {
	return estrategia.Log(msg)
}

type TipoEstrategiaStruct struct {
	estrategy PersistencyStrategyI
}

// Permite definir a estratégia a ser utilizada.
func (e *TipoEstrategiaStruct) SetStrategy(strategy PersistencyStrategyI) {
	e.estrategy = strategy
}

// Permite ler a estratégia utilizada.
func (e *TipoEstrategiaStruct) GetStrategy() PersistencyStrategyI {
	return e.estrategy
}
