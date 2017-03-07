# polymerase

Polymerase is a tool for easy templating using environment variables and vault values.

Polymerase takes a file containing go-style template directives as an argument, populates the template directives with values based on environment variables and vault, and outputs the result to STDOUT.

## Installation

```
go install github.com/dollarshaveclub/polymerase
```

## Usage

```
Usage:
  polymerase [flags]

Examples:
polymerase <filename>

Flags:
  -a, --app-id string         Vault App-ID. Can use APP_ID environment variable instead.
  -u, --user-id-path string   Path to user id. Can use USER_ID_PATH environment variable instead.
  -v, --vault-addr string     Vault server address (including protocol and port). Can use VAULT_ADDR environment variable instead.
  -t, --vault-token string    Vault token. Can use VAULT_TOKEN environment variable instead.
```

## Example


Given `file.tmpl`:

```
Hello, {{ .LOCATION }}! My name is {{ vault "secret_agents/007/first_name" }}.
```

Running:

```
VAULT_ADDR=https://vault.internal VAULT_TOKEN=1234kasd LOCATION=World polymerase file.tmpl
```

Produces:

```
Hello, World! My name is James.
```

Alternatively, input can be provided via stdin:

```
echo "{{ .TEST }}" | VAULT_ADDR=https://vault.internal VAULT_TOKEN=1234kasd LOCATION=World polymerase
```
