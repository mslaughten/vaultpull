// Package sync provides synchronisation between HashiCorp Vault and local
// .env files.
//
// # Schema Validation
//
// SchemaValidator checks a secret map against a declarative set of field
// rules defined in JSON.  Each rule may mark a key as required and/or
// constrain its value with a regular-expression pattern.
//
// Example schema file (.vaultschema.json):
//
//	[
//	  {"key": "DATABASE_URL", "required": true, "pattern": "^postgres://"},
//	  {"key": "PORT",         "required": true, "pattern": "^[0-9]+$"},
//	  {"key": "LOG_LEVEL"}
//	]
//
// Use LoadSchemaFile to read the schema from disk, then call Validate with
// the resolved secret map.  WriteReport prints a human-readable summary to
// any io.Writer.
package sync
