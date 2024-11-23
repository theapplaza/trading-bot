package core

import "trading-bot/common"

type Strategy interface {
	Pass(common.Quote) (float64, bool)
}