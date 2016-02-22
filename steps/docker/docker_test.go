package docker

import (
	"testing"
)

func TestDockerVersion(t *testing.T) {
	opts := &dockerOptions{
		UseTypeNotify:  false,
		DockerExecArgs: make([]string, 0),
	}
	var options *dockerOptions

	options = addVersionSpecificArguments(opts, "1.6.2")
	result := *options
	if result.UseTypeNotify {
		t.Fatalf("expected UseTypeNotify to be false. %v", result.UseTypeNotify)
	}
	if result.DockerExecArgs[0] != "-d" {
		t.Fatalf("expected first argument docker exec to be '-d'. %s", result.DockerExecArgs)
	}

	opts1 := &dockerOptions{
		UseTypeNotify:  true,
		DockerExecArgs: make([]string, 0),
	}
	options = addVersionSpecificArguments(opts1, "1.9.2")
	result = *options
	if !result.UseTypeNotify {
		t.Fatalf("expected UseTypeNotify to be true. %v", result.UseTypeNotify)
	}
	if result.DockerExecArgs[0] != "-d" {
		t.Fatalf("expected first argument docker exec to be '-d'. %s", result.DockerExecArgs[0])
	}
	if result.DockerExecArgs[1] != "--icc=true" {
		t.Fatalf("expected first argument docker exec to be 'daemon'. %s", result.DockerExecArgs[1])
	}

	opts2 := &dockerOptions{
		UseTypeNotify:  true,
		DockerExecArgs: []string{},
	}
	opts2 = addVersionSpecificArguments(opts2, "1.10.2")
	result = *opts2
	if !result.UseTypeNotify {
		t.Fatalf("expected UseTypeNotify to be true. %v", result.UseTypeNotify)
	}
	if result.DockerExecArgs[0] != "daemon" {
		t.Fatalf("expected first argument docker exec to be 'daemon'. %s", result.DockerExecArgs[0])
	}
	if result.DockerExecArgs[1] != "--icc=true" {
		t.Fatalf("expected first argument docker exec to be 'daemon'. %s", result.DockerExecArgs[1])
	}
}
