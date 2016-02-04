package cli

import (
	"github.com/juju/errgo"
)

var (
	ErrWrongInputError = errgo.New("wrong input")

	Mask = errgo.MaskFunc()
)
