package overlay

import (
	"github.com/juju/errgo"
)

var (
	Mask = errgo.MaskFunc(errgo.Any)
)
