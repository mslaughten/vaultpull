// Package sync provides the MaskEnvRenderer which renders a secret map
// to an io.Writer with values partially or fully masked.
//
// Masking modes:
//
//	"all"    – replace every character with the mask symbol (default)
//	"suffix" – reveal the first N characters, mask the rest
//	"prefix" – mask the first characters, reveal the last N characters
//
// Usage:
//
//	r, err := sync.NewMaskEnvRenderer(sync.MaskEnvOptions{
//	    Mode:        sync.MaskEnvSuffix,
//	    RevealChars: 4,
//	    MaskSymbol:  "#",
//	    Keys:        []string{"DB_PASSWORD"},
//	}, os.Stdout)
//	if err != nil { ... }
//	r.Render(secrets)
package sync
