# polymerase

Polymerase is a CLI tool for easy templating using environment variables and [Vault](https://www.vaultproject.io) values.

Polymerase takes a file containing [Go-style template directives `{{ }}`](https://golang.org/pkg/text/template/) as an argument, populates the template directives with values based on environment variables and Vault, and outputs the result to stdout. Input can also be provided via stdin. 

Supported Vault auth backends include [token](https://www.vaultproject.io/docs/auth/token.html) and [App ID](https://www.vaultproject.io/docs/auth/app-id.html). Additionally, [default Go template functions](https://golang.org/pkg/text/template/#hdr-Functions) are supported out of the box. 

<hr />
  <p align="center">
    <a href="#installation">Installation</a>&nbsp;&nbsp;
    <a href="#usage">Usage</a>&nbsp;&nbsp;
    <a href="#examples">Examples</a>&nbsp;&nbsp;
  </p>
<hr />

 
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

## Examples
### File example

Given there is a file with the name `file.tmpl` and contents:

```
Hello, {{ .LOCATION }}! My name is {{ vault "secret_agents/007/first_name" }}.
```

Running the command:

```
VAULT_ADDR=https://vault.internal VAULT_TOKEN=1234kasd LOCATION=World polymerase file.tmpl
```

Polymerase will produce:

```
Hello, World! My name is James.
```

### Stdin example

Running the command:

```
echo "Hello, {{ .LOCATION }}!" | VAULT_ADDR=https://vault.internal VAULT_TOKEN=1234kasd LOCATION=World polymerase
```

Polymerase will produce:

```
Hello, World!
```
