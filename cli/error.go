package cli

import (
	"github.com/juju/errgo"
)

var (
	ErrWrongInputError = errgo.New("wrong input")

	mask = errgo.MaskFunc()
)
