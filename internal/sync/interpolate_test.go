package sync

import (
	"testing"
)

func TestNewInterpolator_ResolvesVariable(t *testing.T) {
	lookup := map[string]string{"HOST": "localhost", "PORT": "5432"}
	ip := NewInterpolator(lookup)

	secrets := map[string]string{
		"DSN": "postgres://${HOST}:${PORT}/db",
	}
	out, err := ip.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["DSN"], "postgres://localhost:5432/db"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolator_DollarSyntax(t *testing.T) {
	lookup := map[string]string{"REGION": "us-east-1"}
	ip := NewInterpolator(lookup)

	out, err := ip.Apply(map[string]string{"BUCKET": "my-bucket-$REGION"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["BUCKET"], "my-bucket-us-east-1"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolator_MissingVar_LenientKeepsOriginal(t *testing.T) {
	ip := NewInterpolator(map[string]string{})

	out, err := ip.Apply(map[string]string{"URL": "http://${HOST}/path"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["URL"], "http://${HOST}/path"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolator_MissingVar_StrictReturnsError(t *testing.T) {
	ip := NewInterpolator(map[string]string{}, WithStrictInterpolation())

	_, err := ip.Apply(map[string]string{"URL": "http://${MISSING}/path"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInterpolator_NoReferences_Unchanged(t *testing.T) {
	ip := NewInterpolator(map[string]string{"X": "y"})

	out, err := ip.Apply(map[string]string{"KEY": "plain-value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["KEY"], "plain-value"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolator_KeysAreUnchanged(t *testing.T) {
	lookup := map[string]string{"V": "x"}
	ip := NewInterpolator(lookup)

	out, err := ip.Apply(map[string]string{"MY_KEY": "${V}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected key MY_KEY to be present")
	}
}
