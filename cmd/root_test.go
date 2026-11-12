package cmd_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tcurvelo/akhq-cli/cmd"
)

// fakeAKHQ serves hand-written fixtures and records each request URI.
type fakeAKHQ struct {
	*httptest.Server
	requests []string
}

func newFakeAKHQ(t *testing.T, routes map[string]string) *fakeAKHQ {
	t.Helper()
	f := &fakeAKHQ{}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.requests = append(f.requests, r.Method+" "+r.RequestURI)
		if r.URL.Path == "/login" {
			http.SetCookie(w, &http.Cookie{Name: "JWT", Value: "test-token"})
			return
		}
		body, ok := routes[r.RequestURI]
		if !ok {
			http.NotFound(w, r)
			return
		}
		if c, err := r.Cookie("JWT"); err != nil || c.Value != "test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(f.Close)
	return f
}

func run(t *testing.T, akhq *fakeAKHQ, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	env := map[string]string{
		"AKHQ_URL":     akhq.URL,
		"AKHQ_CLUSTER": "local",
		"AKHQ_USER":    "alice",
		"AKHQ_PASS":    "secret",
	}
	var out, errOut bytes.Buffer
	err = cmd.NewRootCmd(args, func(k string) string { return env[k] }, &out, &errOut).Execute()
	return out.String(), errOut.String(), err
}

const groupsFixture = `[
  {"id": "billing", "offsets": [
    {"partition": 0, "offsetLag": 3},
    {"partition": 1, "offsetLag": 4}
  ]},
  {"id": "audit", "offsets": [
    {"partition": 0, "offsetLag": 10},
    {"partition": 1, "offsetLag": null}
  ]}
]`

func TestGroupsPrintsOffsetLagPerGroup(t *testing.T) {
	akhq := newFakeAKHQ(t, map[string]string{
		"/api/local/topic/orders/groups": groupsFixture,
	})

	stdout, _, err := run(t, akhq, "groups", "orders")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "7\n10\n"; stdout != want {
		t.Errorf("stdout = %q, want %q", stdout, want)
	}
}

func TestGroupsEscapesTopicName(t *testing.T) {
	akhq := newFakeAKHQ(t, map[string]string{
		"/api/local/topic/orders%2Fdlq%20%231/groups": groupsFixture,
	})

	stdout, _, err := run(t, akhq, "groups", "orders/dlq #1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "7\n10\n"; stdout != want {
		t.Errorf("stdout = %q, want %q (requests: %v)", stdout, want, akhq.requests)
	}
}
