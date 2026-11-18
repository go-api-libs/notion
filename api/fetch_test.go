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

// TestLikelyNow tells, for a name in shortNames or exactNames that flatten no longer gives, what it likely gives
// instead: the name ending in most of the same words, as when upstream titled the response holding it.
func TestLikelyNow(t *testing.T) {
	t.Parallel()

	names := []string{
		"AuthorizationCodeOrRefreshTokenOwnerUserUser",
		"AuthorizationCodeOrRefreshTokenOwnerUser",
		"BotInfoResponseOwnerUserUser2",
		"PageComment",
	}

	for old, want := range map[string]string{
		"CreateATokenOkOwnerUserUser": "AuthorizationCodeOrRefreshTokenOwnerUserUser",
		"CreateATokenOkOwnerUser":     "AuthorizationCodeOrRefreshTokenOwnerUser",
		"CreateACommentAllOf0OneOf0":  "", // no two words alike
	} {
		if got := likelyNow(old, names); got != want {
			t.Errorf("likelyNow(%q) = %q, want %q", old, got, want)
		}
	}

	if got := likelyNow("XOwnerUser", []string{"AOwnerUser", "BOwnerUser"}); got != "" {
		t.Errorf("of two alike: %q, want none", got)
	}
}
