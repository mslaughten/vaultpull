package sync

import (
	"errors"
	"testing"
)

func TestResult_HasErrors(t *testing.T) {
	r := Result{}
	if r.HasErrors() {
		t.Error("empty result should not have errors")
	}

	r.Errors = append(r.Errors, errors.New("oops"))
	if !r.HasErrors() {
		t.Error("result with errors should return true")
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Total:   3,
		Written: []string{"a.env", "b.env"},
		Errors:  []error{errors.New("fail")},
	}
	got := r.Summary()
	want := "synced 2/3 secret(s), 1 error(s)"
	if got != want {
		t.Errorf("Summary() = %q; want %q", got, want)
	}
}

func TestResult_ErrorMessages(t *testing.T) {
	r := Result{
		Errors: []error{errors.New("err1"), errors.New("err2")},
	}
	got := r.ErrorMessages()
	if got != "err1\nerr2" {
		t.Errorf("ErrorMessages() = %q", got)
	}
}

func TestEnvFileName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"db", "db.env"},
		{"app/prod", "app_prod.env"},
		{"", "secrets.env"},
	}
	for _, tc := range cases {
		got := envFileName(tc.input)
		if got != tc.want {
			t.Errorf("envFileName(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}
