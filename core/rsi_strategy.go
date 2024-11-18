package core

type RsiStrategy struct {
    Name   string
    Period int
    Level  float64
}


func (rsi *RsiStrategy) Pass() bool {
	return true
}