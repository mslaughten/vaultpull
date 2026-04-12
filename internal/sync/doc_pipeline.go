// Package sync provides synchronisation primitives for vaultpull.
//
// # Pipeline
//
// A Pipeline chains named transformation stages that each receive a
// map[string]string of secrets and return a (possibly modified) copy.
//
// Stages are applied in the order they were added. If any stage returns an
// error the pipeline halts immediately and returns that error annotated with
// the failing stage's name.
//
// Usage:
//
//	p := sync.NewPipeline()
//	p.AddStage(sync.PipelineStage{Name: "upper", Apply: upperFunc})
//	p.AddStage(sync.PipelineStage{Name: "redact", Apply: redactFunc})
//	out, err := p.Run(secrets)
package sync
