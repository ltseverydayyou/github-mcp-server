package github

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSandboxPath(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GITHUB_MCP_SANDBOX_ROOT", root)
	path := filepath.Join(root, "large.lua")
	if err := os.WriteFile(path, []byte("print('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}
	resolved, err := resolveSandboxPath(path, true)
	if err != nil {
		t.Fatalf("resolveSandboxPath returned error: %v", err)
	}
	if resolved != path {
		t.Fatalf("unexpected resolved path: %q", resolved)
	}
	outside := filepath.Join(filepath.Dir(root), "outside.lua")
	if _, err := resolveSandboxPath(outside, false); err == nil {
		t.Fatal("expected outside path to be rejected")
	}
}

func TestResolveSandboxDestination(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GITHUB_MCP_SANDBOX_ROOT", root)
	destination := filepath.Join(root, "nested", "source.go")
	resolved, err := resolveSandboxDestination(destination)
	if err != nil {
		t.Fatalf("resolveSandboxDestination returned error: %v", err)
	}
	if resolved != destination {
		t.Fatalf("unexpected resolved destination: %q", resolved)
	}
}
