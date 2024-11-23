package core

type MacdStrategy struct {
	Name string
}

func (rsi *MacdStrategy) Check() (float64, bool) {
	return 0, false
}