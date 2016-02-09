package fs

import (
	"os"

	"github.com/juju/errgo"
)

var (
	mask = errgo.MaskFunc()
)

func IsNotExist(err error) bool {
	return os.IsNotExist(errgo.Cause(err))
}

func IsExist(err error) bool {
	return os.IsExist(errgo.Cause(err))
}
