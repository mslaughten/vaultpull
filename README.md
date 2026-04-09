# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files with namespace filtering

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Set your Vault address and token, then run `vaultpull` with a namespace path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a namespace into a local .env file
vaultpull --namespace secret/myapp/production --output .env
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--namespace` | Vault secret path / namespace to pull from | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--filter` | Comma-separated list of key prefixes to include | *(all keys)* |
| `--overwrite` | Overwrite existing `.env` file | `false` |

**Example output (`.env`):**
```
DB_HOST=prod-db.internal
DB_PASSWORD=supersecret
API_KEY=abc123
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance (v1.x)
- A valid Vault token with read access to the target namespace

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)