package core

import "trading-bot/common"

type Strategy interface {
	Check(common.Quote) (float64, bool)
}