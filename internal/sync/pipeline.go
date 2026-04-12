package sync

import (
	"fmt"
	"strings"
)

// PipelineStage represents a named transformation step applied to a secret map.
type PipelineStage struct {
	Name string
	Apply func(map[string]string) (map[string]string, error)
}

// Pipeline chains multiple stages together, applying each in order.
type Pipeline struct {
	stages []PipelineStage
}

// NewPipeline constructs an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage appends a stage to the pipeline.
func (p *Pipeline) AddStage(stage PipelineStage) *Pipeline {
	if strings.TrimSpace(stage.Name) == "" {
		stage.Name = fmt.Sprintf("stage-%d", len(p.stages)+1)
	}
	p.stages = append(p.stages, stage)
	return p
}

// Run executes all stages in order, passing the output of each stage as the
// input of the next. It returns the final map or the first error encountered,
// annotated with the failing stage name.
func (p *Pipeline) Run(input map[string]string) (map[string]string, error) {
	current := copyMap(input)
	for _, stage := range p.stages {
		out, err := stage.Apply(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", stage.Name, err)
		}
		current = out
	}
	return current, nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int { return len(p.stages) }

// StageNames returns the ordered list of stage names.
func (p *Pipeline) StageNames() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
