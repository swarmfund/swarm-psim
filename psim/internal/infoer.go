package internal

import "gitlab.com/tokend/regources"

//go:generate mockery -case underscore -name Infoer

// Infoer is capable of providing TokenD network info
type Infoer interface {
	Info() (*regources.Info, error)
}
