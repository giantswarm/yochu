package templates

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	// We picked the usr-bin template since it is rather simple and contains only a single value: "ExecStart"
	asset := "templates/usr-bin.mount.tmpl"
	data := map[string]string{
		"MountPoint": "THIS_IS_A_TEST",
	}
	content, err := Render(asset, data)
	if err != nil {
		t.Fatalf("rendering of %s failed: %v", asset, err.Error())
	}

	if !strings.Contains(content.String(), data["ExecStart"]) {
		t.Fatalf("expected rendered template to contain passed value.")
	}
}
