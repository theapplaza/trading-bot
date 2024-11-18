package core

type MacdStrategy struct {
	Name string
}

func (rsi *MacdStrategy) Pass() bool {
	return true
}