// Package templates provides functions to manage templates.
package templates

import (
	"bytes"
	"text/template"

	"github.com/juju/errgo"
)

var maskAny = errgo.MaskFunc(errgo.Any)

// Render takes an asset name of a template and some arguments,
// and renders and returns the template.
func Render(assetName string, arguments interface{}) (*bytes.Buffer, error) {
	templateRaw, err := Asset(assetName)
	if err != nil {
		return nil, maskAny(err)
	}

	var tmpl *template.Template
	tmpl, err = template.New("template").Parse(string(templateRaw))
	if err != nil {
		return nil, maskAny(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, arguments)
	if err != nil {
		return nil, maskAny(err)
	}
	return &b, nil
}
