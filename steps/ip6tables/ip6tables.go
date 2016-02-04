package ip6tables

import (
	"os"

	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

type ip6tablesOptions struct {
}

var (
	vLogger = func(f string, v ...interface{}) {}

	netfilterName = "netfilter.service"
	serviceName   = "ip6tables.service"
	rulesName     = "ip6tables.rules"

	netfilterTemplate = "templates/netfilter.service.tmpl"
	serviceTemplate   = "templates/ip6tables.service.tmpl"
	rulesTemplate     = "templates/ip6tables.rules.tmpl"

	netfilterPath = "/etc/systemd/system/" + netfilterName
	servicePath   = "/etc/systemd/system/" + serviceName
	rulesPath     = "/home/core/" + rulesName

	fileMode = os.FileMode(0755)
	services = []string{serviceName}
	paths    = []string{servicePath, rulesPath}
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsClient *fs.FsClient, systemdClient *systemd.SystemdClient) error {
	vLogger("\n# call ip6tables.Setup()")

	if err := setupRules(fsClient); err != nil {
		return Mask(err)
	}

	if err := setupService(fsClient); err != nil {
		return Mask(err)
	}

	if err := systemdClient.Reload(); err != nil {
		return Mask(err)
	}

	for _, s := range services {
		if err := systemdClient.Start(s); err != nil {
			return Mask(err)
		}
	}

	return nil
}

func RenderRulesFromTemplate() ([]byte, error) {
	opts := ip6tablesOptions{}
	b, err := templates.Render(rulesTemplate, opts)
	if err != nil {
		return nil, Mask(err)
	}
	return b.Bytes(), nil
}

func Teardown(fsClient *fs.FsClient, systemdClient *systemd.SystemdClient) error {
	vLogger("\n# call ip6tables.Teardown()")

	for _, s := range services {
		exists, err := systemdClient.Exists(s)
		if err != nil {
			return Mask(err)
		}

		if !exists {
			continue
		}

		if err := systemdClient.Stop(s); err != nil {
			return Mask(err)
		}
	}

	for _, p := range paths {
		if err := fsClient.Remove(p); err != nil {
			return Mask(err)
		}
	}

	if err := systemdClient.Reload(); err != nil {
		return Mask(err)
	}

	return nil
}

func setupRules(fsClient *fs.FsClient) error {
	rulesBytes, err := RenderRulesFromTemplate()
	if err != nil {
		return Mask(err)
	}

	if err := fsClient.Write(rulesPath, rulesBytes, fileMode); err != nil {
		return Mask(err)
	}

	return nil
}

func setupService(fsClient *fs.FsClient) error {
	iptablesService, err := templates.Asset(serviceTemplate)
	if err != nil {
		return Mask(err)
	}

	if err := fsClient.Write(servicePath, iptablesService, fileMode); err != nil {
		return Mask(err)
	}

	return nil
}
