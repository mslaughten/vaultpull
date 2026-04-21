// Package sync provides the LabelMapper type for renaming secret keys
// using a static label-to-key mapping.
//
// # Overview
//
// LabelMapper accepts a list of "label=newkey" rules. When applied to a
// map of secrets, any key that matches a label is emitted under the
// corresponding new key name. Keys without a matching rule are passed
// through unchanged.
//
// # Usage
//
//	lm, err := sync.NewLabelMapper([]string{
//	    "DB_PASS=DATABASE_PASSWORD",
//	    "API=API_KEY",
//	})
//	result, err := lm.Apply(secrets)
package sync
