package sync

import (
	"testing"
)

func TestInvertStrategyFromString_Valid(t *testing.T) {
 _, s := range []string{"keys", "values",t	g	if err != nil {
			t.Fatalf("unexpected error for %q: %v", s, err)
		}
		if got != s {
			t.Errorf("expected %q, got %q", s, got)
		}
	}
}

func TestInvertStrategyFromString_Invalid(t *testing.T) {
	_, err := InvertStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewInverter_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewInverter("flip")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInverter_Keys_ReversesKeys(t *testing.T) {
	inv, err := NewInverter("keys")
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{"ABC": "hello", "XY": "world"}
	out, err := inv.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out["CBA"] != "hello" {
		t.Errorf("expected key CBA=hello, got %v", out)
	}
	if out["YX"] != "world" {
		t.Errorf("expected key YX=world, got %v", out)
	}
}

func TestInverter_Values_ReversesValues(t *testing.T) {
	inv, err := NewInverter("values")
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{"KEY": "abc"}
	out, err := inv.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "cba" {
		t.Errorf("expected cba, got %q", out["KEY"])
	}
}

func TestInverter_Both_ReversesBothKeyAndValue(t *testing.T) {
	inv, err := NewInverter("both")
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{"ENV": "prod"}
	out, err := inv.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if out["VNE"] != "dorp" {
		t.Errorf("expected VNE=dorp, got %v", out)
	}
}

func TestInverter_EmptyMap_ReturnsEmpty(t *testing.T) {
	inv, err := NewInverter("both")
	if err != nil {
		t.Fatal(err)
	}
	out, err := inv.Apply(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
