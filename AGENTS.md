# gotest-labels

A Go library for filtering Go test cases by labels in test comments.

## Project Rules

- **Keep dependencies small** — the current external dependency is `golang.org/x/tools` for Go package loading; avoid adding more unless clearly justified.
- **Keep README examples accurate** — examples must use `github.com/maxwu/gotest-labels` and the current `MutateTestFilterByLabels()` API.
- **Verify changes** — run `go test -v -count=1 ./...` and `golangci-lint run --timeout 1m` after code changes.
