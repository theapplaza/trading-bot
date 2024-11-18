package core

type Strategy interface {
	Pass() bool
}