package main

import (
	cliPkg "github.com/giantswarm/yochu/cli"
)

var projectVersion string

func main() {
	cliPkg.NewYochuCmd(projectVersion)
}
