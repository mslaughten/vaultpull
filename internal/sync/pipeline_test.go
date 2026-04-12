package sync

import (
	"errors"
	"strings"
	"testing"
)

func TestNewPipeline_EmptyHasZeroStages(t *testing.T) {
	p := NewPipeline()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
}

func TestPipeline_AddStage_IncrementsLen(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{Name: "upper", Apply: func(m map[string]string) (map[string]string, error) { return m, nil }})
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestPipeline_StageNames(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{Name: "a", Apply: func(m map[string]string) (map[string]string, error) { return m, nil }})
	p.AddStage(PipelineStage{Name: "b", Apply: func(m map[string]string) (map[string]string, error) { return m, nil }})
	names := p.StageNames()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected stage names: %v", names)
	}
}

func TestPipeline_Run_AppliesStagesInOrder(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{
		Name: "uppercase",
		Apply: func(m map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(m))
			for k, v := range m {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		},
	})
	p.AddStage(PipelineStage{
		Name: "prefix",
		Apply: func(m map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(m))
			for k, v := range m {
				out[k] = "PRE_" + v
			}
			return out, nil
		},
	})

	result, err := p.Run(map[string]string{"key": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "PRE_HELLO" {
		t.Fatalf("expected PRE_HELLO, got %s", result["key"])
	}
}

func TestPipeline_Run_StageError_AnnotatesName(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{
		Name: "fail",
		Apply: func(m map[string]string) (map[string]string, error) {
			return nil, errors.New("boom")
		},
	})
	_, err := p.Run(map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Fatalf("error should contain stage name, got: %v", err)
	}
}

func TestPipeline_AddStage_EmptyName_AutoAssigned(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{Apply: func(m map[string]string) (map[string]string, error) { return m, nil }})
	names := p.StageNames()
	if names[0] == "" {
		t.Fatal("expected auto-assigned name, got empty string")
	}
}

func TestPipeline_Run_DoesNotMutateInput(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PipelineStage{
		Name: "del",
		Apply: func(m map[string]string) (map[string]string, error) {
			delete(m, "secret")
			return m, nil
		},
	})
	input := map[string]string{"secret": "value"}
	p.Run(input) //nolint:errcheck
	if _, ok := input["secret"]; !ok {
		t.Fatal("Run mutated the input map")
	}
}
