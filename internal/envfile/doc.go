// Package envfile provides functionality for writing secrets to .env files.
//
// The package handles:
//   - Creating .env files with proper permissions (0600)
//   - Sorting keys alphabetically for consistent output
//   - Escaping special characters in values (quotes, newlines, backslashes)
//   - Creating parent directories if they don't exist
//   - Adding header comments to identify generated files
//
// Example usage:
//
//	writer := envfile.NewWriter(".env")
//	secrets := map[string]string{
//		"DATABASE_URL": "postgres://localhost/db",
//		"API_KEY": "secret123",
//	}
//	if err := writer.Write(secrets); err != nil {
//		log.Fatal(err)
//	}
//
// The output .env file will have restricted permissions (0600) to protect
// sensitive secrets from unauthorized access.
package envfile
