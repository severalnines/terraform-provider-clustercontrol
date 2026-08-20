# Makefile targets — `terraform-provider-clustercontrol`

Location: repo root (`./Makefile`)

Run `make help` at any time to print this list from the Makefile itself.

| Target | Description |
|---|---|
| `build` | Compile the provider binary for the current OS/ARCH into `./bin/terraform-provider-clustercontrol`. Uses `CGO_ENABLED=0`. |
| `install` *(default)* | Runs `build`, then copies the binary into the local Terraform dev-override plugin directory: `~/.terraform.d/plugins/severalnines.com/severalnines/clustercontrol/<VERSION>/<OS>_<ARCH>/`. `<OS>_<ARCH>` is auto-detected via `go env` (e.g. `linux_amd64` on Ubuntu) rather than hardcoded. |
| `clean` | Removes the installed dev-override binary directory and the local `./bin` build output. |
| `fmt` | Runs `go fmt ./...` across the module. |
| `vet` | Runs `go vet ./...` across the module. |
| `lint` | Runs `golangci-lint run ./...` if `golangci-lint` is on `PATH`; otherwise prints an install hint and exits cleanly (non-fatal). |
| `tidy` | Runs `go mod tidy`. |
| `test` | Runs unit tests (`go test ./... -v -count=1`). No live ClusterControl instance required. |
| `testacc` | Runs acceptance tests with `TF_ACC=1` and a longer timeout (`-timeout 120m`). Requires a reachable ClusterControl instance and valid `cc_api_url`/`cc_api_user`/`cc_api_user_password` credentials configured however your acceptance tests read them (env vars/tfvars). |
| `docs` | Regenerates `./docs` from resource schema descriptions and `./examples`, via `tfplugindocs` (`go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs`). |
| `release` | Runs a local snapshot release with `goreleaser` (`--clean --snapshot --skip=publish --skip=sign`) — builds all cross-platform binaries without publishing or signing. |
| `sdk-update` | Bumps the `clustercontrol-client-sdk` Go dependency to latest (`go get -u .../go/pkg/openapi`) and runs `go mod tidy`. Use this after the SDK repo publishes a new tag (e.g. after adding ClickHouse support upstream). |
| `help` | Prints all `##`-annotated targets plus the detected `OS_ARCH`. |

## Variables you can override

| Variable | Default | Purpose |
|---|---|---|
| `VERSION` | `0.2.23` | Version segment of the dev-override install path. Override to test a different version string, e.g. `make install VERSION=0.3.0`. |
| `GO` | `go` | Path to the Go toolchain. |
| `GOLANGCI_LINT` | `golangci-lint` | Path to the linter binary. |

## Typical workflows

```sh
# Local dev loop
make fmt vet lint test
make install                # builds + installs into dev-override path
terraform init && terraform plan   # picks up the dev-override provider

# Refresh generated docs after changing schema Description fields
make docs

# Pull in a newer SDK (e.g. after ClickHouse support lands upstream)
make sdk-update
make build
```
