// Package sync provides the core synchronisation logic for vaultpull.
//
// # Interpolation
//
// The Interpolator type performs variable substitution inside secret values
// before they are written to .env files. Both ${VAR} and $VAR syntaxes are
// supported.
//
// References are resolved against a caller-supplied lookup map, which is
// typically the full set of secrets already fetched from Vault, or a
// snapshot of the process environment.
//
// By default missing variables are left as-is (lenient mode). Pass
// WithStrictInterpolation to treat any unresolved reference as an error.
//
// # Strict Mode
//
// When strict mode is enabled, Apply returns an [InterpolationError] that
// lists every unresolved variable name encountered during the substitution
// pass. This makes it easy to surface configuration mistakes early rather
// than silently writing incomplete values to disk.
//
// Example:
//
//		lookup := map[string]string{"HOST": "db.internal", "PORT": "5432"}
//		ip := sync.NewInterpolator(lookup)
//		out, err := ip.Apply(secrets)
//
//		// strict mode – any unresolved reference is an error
//		ip2 := sync.NewInterpolator(lookup, sync.WithStrictInterpolation())
//		out, err = ip2.Apply(secrets)
package sync
