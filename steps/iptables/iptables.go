package iptables

import (
	"os"

	"github.com/giantswarm/yochu/fs"
	"github.com/giantswarm/yochu/systemd"
	"github.com/giantswarm/yochu/templates"
)

type iptablesOptions struct {
	Gateway        string
	Subnet         string
	DockerSubnet   string
	UseDockerRules bool
}

var (
	vLogger = func(f string, v ...interface{}) {}

	netfilterName = "netfilter.service"
	serviceName   = "iptables.service"
	rulesName     = "iptables.rules"

	netfilterTemplate = "templates/netfilter.service.tmpl"
	serviceTemplate   = "templates/iptables.service.tmpl"
	rulesTemplate     = "templates/iptables.rules.tmpl"

	netfilterPath = "/etc/systemd/system/" + netfilterName
	servicePath   = "/etc/systemd/system/" + serviceName
	rulesPath     = "/home/core/" + rulesName

	fileMode = os.FileMode(0755)
	services = []string{netfilterName, serviceName}
	paths    = []string{netfilterPath, servicePath, rulesPath}
)

func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

func Setup(fsc *fs.FsClient, sc *systemd.SystemdClient, subnet, dockerSubnet, gateway string, useDockerRules bool) error {
	vLogger("\n# call iptables.Setup()")

	if err := setupNetfilter(fsc); err != nil {
		return maskAny(err)
	}

	if err := setupRules(fsc, subnet, dockerSubnet, gateway, useDockerRules); err != nil {
		return maskAny(err)
	}

	if err := setupService(fsc); err != nil {
		return maskAny(err)
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	for _, s := range services {
		if err := sc.Start(s); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func RenderRulesFromTemplate(subnet, dockerSubnet, gateway string, useDockerRules bool) ([]byte, error) {
	opts := iptablesOptions{
		Subnet:         subnet,
		DockerSubnet:   dockerSubnet,
		Gateway:        gateway,
		UseDockerRules: useDockerRules,
	}

	b, err := templates.Render(rulesTemplate, opts)
	if err != nil {
		return nil, maskAny(err)
	}
	return b.Bytes(), nil
}

func Teardown(fsc *fs.FsClient, sc *systemd.SystemdClient) error {
	vLogger("\n# call iptables.Teardown()")

	for _, s := range services {
		exists, err := sc.Exists(s)
		if err != nil {
			return maskAny(err)
		}

		if !exists {
			continue
		}

		if err := sc.Stop(s); err != nil {
			return maskAny(err)
		}
	}

	for _, p := range paths {
		if err := fsc.Remove(p); err != nil {
			return maskAny(err)
		}
	}

	if err := sc.Reload(); err != nil {
		return maskAny(err)
	}

	return nil
}

func setupNetfilter(fsc *fs.FsClient) error {
	netfilterService, err := templates.Asset(netfilterTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(netfilterPath, netfilterService, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func setupRules(fsc *fs.FsClient, subnet, dockerSubnet, gateway string, useDockerRules bool) error {
	rulesBytes, err := RenderRulesFromTemplate(subnet, dockerSubnet, gateway, useDockerRules)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(rulesPath, rulesBytes, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}

func setupService(fsc *fs.FsClient) error {
	iptablesService, err := templates.Asset(serviceTemplate)
	if err != nil {
		return maskAny(err)
	}

	if err := fsc.Write(servicePath, iptablesService, fileMode); err != nil {
		return maskAny(err)
	}

	return nil
}
