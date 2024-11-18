package core

import "trading-bot/common"

type Strategy interface {
	Pass(producer string, symbol common.Symbol) (float64, bool)
}