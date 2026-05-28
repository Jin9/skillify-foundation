# Make Commands

Repository targets defined in the root `Makefile`. Run from the `common` repo root.

| Command | Purpose |
|---|---|
| `make setup` | Install dev tools and git hooks |
| `make mod` | `go fmt` + `go mod tidy` |
| `make lint` | `golangci-lint` with auto-fix |
| `make test` | Tests with `-race`, coverage, colorised output |
| `make coverage` | Generate HTML coverage report |
| `make vuln` | `govulncheck` for known vulnerabilities |
| `make precommit` | Full validation: lint + test + vuln + vet + mod verify |
| `make ci` | CI pipeline: precommit + uncommitted-changes check |
| `make doc` | Serve package docs via `pkgsite` on `:6060` |
| `make upgrade` | Upgrade dev tools |
| `make bump-version version=vX.Y.Z` | Tag a new release |

Always run `make precommit` before pushing.
