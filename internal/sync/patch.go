package sync

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation on a secret map.
type PatchOp struct {
	Key   string
	Op    string // set, delete, append
	Value string
}

// Patcher applies a list of patch operations to a secret map.
type Patcher struct {
	ops []PatchOp
}

// NewPatcher creates a Patcher from a slice of operation strings.
// Each string must have the format "op:key=value" or "delete:key".
func NewPatcher(rules []string) (*Patcher, error) {
	ops := make([]PatchOp, 0, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("patch: invalid rule %q: expected op:key or op:key=value", r)
		}
		op := strings.TrimSpace(parts[0])
		rest := parts[1]
		switch op {
		case "set", "append":
			kv := strings.SplitN(rest, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("patch: rule %q requires key=value for op %q", r, op)
			}
			key := strings.TrimSpace(kv[0])
			if key == "" {
				return nil, fmt.Errorf("patch: empty key in rule %q", r)
			}
			ops = append(ops, PatchOp{Key: key, Op: op, Value: kv[1]})
		case "delete":
			key := strings.TrimSpace(rest)
			if key == "" {
				return nil, fmt.Errorf("patch: empty key in rule %q", r)
			}
			ops = append(ops, PatchOp{Key: key, Op: op})
		default:
			return nil, fmt.Errorf("patch: unknown op %q in rule %q", op, r)
		}
	}
	return &Patcher{ops: ops}, nil
}

// Apply runs all patch operations against the provided map and returns
// a new map with the changes applied.
func (p *Patcher) Apply(in map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	for _, op := range p.ops {
		switch op.Op {
		case "set":
			out[op.Key] = op.Value
		case "delete":
			delete(out, op.Key)
		case "append":
			if existing, ok := out[op.Key]; ok {
				out[op.Key] = existing + op.Value
			} else {
				out[op.Key] = op.Value
			}
		}
	}
	return out, nil
}
