package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerate makes the specification and the library from the official specification and the recordings
// committed, offline, as go generate does, and finds them as committed: what the repository holds is what the
// generator makes of what it holds, no file of it edited by hand or left behind.
func TestGenerate(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	tmp := t.TempDir()

	for _, name := range []string{"api/openapi-official.json", "api/interactions.json"} {
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll(filepath.Join(tmp, filepath.Dir(name)), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(tmp, name), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Chdir(tmp)

	isClaudeCode = true // read the official specification committed rather than fetch it

	if err := persistOpenAPI(t.Context()); err != nil {
		t.Fatal(err)
	}

	doc, err := fixOpenAPI()
	if err != nil {
		t.Fatal(err)
	}

	if err := generateCode(doc); err != nil {
		t.Fatal(err)
	}

	generated, err := filepath.Glob("pkg/notion/*")
	if err != nil {
		t.Fatal(err)
	}

	if len(generated) == 0 {
		t.Fatal("nothing generated")
	}

	for _, name := range append(generated, pathOpenAPI) {
		got, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		want, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Errorf("%s is generated but not committed: %v", name, err)
			continue
		}

		if !bytes.Equal(got, want) {
			t.Errorf("%s differs from what the generator makes: run go generate", name)
		}
	}
}
