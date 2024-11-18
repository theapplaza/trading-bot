package core

type MacdStrategy struct {
	Name string
}

func (rsi *MacdStrategy) Pass() (float64, bool) {
	return 0, false
}