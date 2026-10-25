package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// run executes fn with stdout captured, against a server that records the
// last request and replies with body.
func run(t *testing.T, body string, fn func(base string) error) (out string, method, path, reqBody string, err error) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		method, path, reqBody = r.Method, r.URL.RequestURI(), string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err = fn(srv.URL)
	w.Close()
	os.Stdout = old
	data, _ := io.ReadAll(r)
	return string(data), method, path, reqBody, err
}

func common(base string) []string { return []string{"--api-key", "dbk_test", "--base-url", base} }

func TestPurgeFile(t *testing.T) {
	_, m, p, b, err := run(t, `{"success":true,"purged":1,"failed":[]}`, func(base string) error {
		return Purge(append(common(base), "--scope", "file", "--key", "a.png", "media"))
	})
	if err != nil || m != "POST" || p != "/purge" {
		t.Fatalf("got %v %s %s", err, m, p)
	}
	if !strings.Contains(b, `"scope":"file"`) || !strings.Contains(b, `"path":"a.png"`) || !strings.Contains(b, `"slug":"media"`) {
		t.Fatalf("body %s", b)
	}
}

func TestPurgeValidation(t *testing.T) {
	for name, args := range map[string][]string{
		"file without key": {"--scope", "file", "media"},
		"bad scope":        {"--scope", "all", "media"},
		"key on bucket":    {"--key", "a.png", "media"},
	} {
		if err := Purge(append([]string{"--api-key", "k"}, args...)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestBucketsUpdateSendsOnlySetFlags(t *testing.T) {
	_, m, p, b, err := run(t, `{"bucket":{}}`, func(base string) error {
		return Buckets(append([]string{"update"}, append(common(base), "--cdn=false", "media")...))
	})
	if err != nil || m != "PATCH" || p != "/buckets/media" || b != `{"cdn":false}` {
		t.Fatalf("got %v %s %s %s", err, m, p, b)
	}
	if err := Buckets([]string{"update", "--api-key", "k", "media"}); err == nil {
		t.Fatal("expected error when nothing to update")
	}
}

func TestBucketsCreateAndDelete(t *testing.T) {
	_, m, _, b, err := run(t, `{"bucket":{}}`, func(base string) error {
		return Buckets(append([]string{"create"}, append(common(base), "--region", "us-east", "--public", "Media")...))
	})
	if err != nil || m != "POST" || !strings.Contains(b, `"isPublic":true`) || !strings.Contains(b, `"name":"Media"`) {
		t.Fatalf("got %v %s %s", err, m, b)
	}
	if err := Buckets([]string{"create", "--api-key", "k", "Media"}); err == nil {
		t.Fatal("expected error without --region")
	}
	if err := Buckets([]string{"delete", "--api-key", "k", "media"}); err == nil {
		t.Fatal("delete without --yes must fail")
	}
	_, m, p, _, err := run(t, `{"success":true}`, func(base string) error {
		return Buckets(append([]string{"delete"}, append(common(base), "--yes", "media")...))
	})
	if err != nil || m != "DELETE" || p != "/buckets/media" {
		t.Fatalf("got %v %s %s", err, m, p)
	}
}

func TestUploadURL(t *testing.T) {
	_, m, p, b, err := run(t, `{"url":"u","method":"PUT","headers":{},"expiresIn":900}`, func(base string) error {
		return UploadURL(append(common(base), "--content-type", "image/png", "--size", "10", "media", "a.png"))
	})
	if err != nil || m != "POST" || p != "/storage/media/upload-url" || !strings.Contains(b, `"name":"a.png"`) {
		t.Fatalf("got %v %s %s %s", err, m, p, b)
	}
}

func TestDataCommands(t *testing.T) {
	_, m, p, _, err := run(t, `{"data":[],"next_cursor":null}`, func(base string) error {
		return Data(append([]string{"get"}, append(common(base), "--limit", "5", "--filter", "city=Oslo", "people")...))
	})
	if err != nil || m != "GET" || !strings.HasPrefix(p, "/v1/data/people?") || !strings.Contains(p, "city=Oslo") || !strings.Contains(p, "limit=5") {
		t.Fatalf("got %v %s %s", err, m, p)
	}
	_, m, p, b, err := run(t, `{"data":{}}`, func(base string) error {
		return Data(append([]string{"update"}, append(common(base), "--json", `{"a":1}`, "people", "7")...))
	})
	if err != nil || m != "PATCH" || p != "/v1/data/people/7" || b != `{"a":1}` {
		t.Fatalf("got %v %s %s %s", err, m, p, b)
	}
	_, m, _, b, err = run(t, `{"success":true,"appendedCharacters":2}`, func(base string) error {
		return Data(append([]string{"append-text"}, append(common(base), "--text", "hi", "notes")...))
	})
	if err != nil || m != "POST" || b != `{"appendText":"hi"}` {
		t.Fatalf("got %v %s %s", err, m, b)
	}
	if err := Data([]string{"create", "--api-key", "k", "--json", "[1]", "people"}); err == nil {
		t.Fatal("non-object JSON must be rejected")
	}
	if err := Data([]string{"create", "--api-key", "k", "people"}); err == nil {
		t.Fatal("missing --json/--file must be rejected")
	}
	if err := Data([]string{"get", "--api-key", "k", "--filter", "bad", "people"}); err == nil {
		t.Fatal("bad filter must be rejected")
	}
}
