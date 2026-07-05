# Exploration: gotest-labels

**Date**: 2026-06-30 | **Scope**: Medium | **Status**: ✅ Complete

## 1. Foundation (What exists)

**Tech stack**: Go 1.26, single module `github.com/maxwu/gotest-labels`, one external dependency `golang.org/x/tools` (for `go/packages`)

**Architecture**: Go library package — not a CLI binary. Provides test-label filtering by mutating `os.Args` before `go test` runs. The filtering injects a `-test.run` or `-test.list` regex pattern into `os.Args` at init time.

**Structure**:
- Root package `gotest_labels` — all core logic (4 source files, 554 LOC)
- `apply/` — thin `init()` wrapper for anonymous import (8 LOC)
- `examples/` — 3 example packages demonstrating usage patterns
- `.history/` — local edit history (not in git)

| File | Purpose | LOC |
|------|---------|-----|
| `gotestlabels.go` | Public entrypoint `MutateTestFilterByLabels()`, orchestrates the full flow | 71 |
| `exp_parser.go` | Expression tokenizer + recursive-descent parser → AST, and `Evaluate()` | 216 |
| `pkg_parser.go` | Go AST parsing: discovers test funcs, extracts `@key=value` labels from comments | 142 |
| `cli.go` | `os.Args` parsing for `-labels`, `-test.run`, `-test.list`; removes custom `-labels` flag | 125 |
| `apply/apply.go` | `init()` auto-apply via anonymous import | 8 |

**CLAUDE.md instructions**: None found (no project-root or `.claude/` CLAUDE.md files).

## 2. Patterns (How it's built)

**Core data flow**:
```
os.Args / TEST_LABELS env
  → cli.go: ParseOSArgs() → cliArgs{labels, runRegex, listMode, labelsAST}
  → exp_parser.go: ParseLabelExp() → AST (Condition | LogicalOp{AND,OR,NOT})
  → pkg_parser.go: getPackages() → []*packages.Package
  → pkg_parser.go: FindTestFuncs(files, filterAST) → map[string]TestLabels
  → pkg_parser.go: filterTestFuncs(matched, runRegex) → map[string]TestLabels
  → gotestlabels.go: MutateTestFilterByLabels() → mutates os.Args with regex pattern
```

**Expression parser**: Recursive-descent with 3 precedence levels — `parseExpr` (OR) > `parseTerm` (AND) > `parseFactor` (NOT/parens/condition). AST nodes: `Condition{Key, Value}` and `LogicalOp{Operator, Children}`.

**Label extraction**: Uses `go/ast` parser on `*_test.go` files. Reads `fn.Doc.List` comments, strips `//` or `/* */`, matches `@key=value` or `@key` (defaults to `true`). Validates `func Test*(t *testing.T)` signatures.

**Testing patterns**:
- Standard `testing` package, table-driven tests with `map[string]struct{}` or `[]struct{}`
- `t.Parallel()` used extensively
- `t.Setenv()` for env-var tests (Go 1.17+)
- Direct `os.Args` mutation with defer-restore for integration-style tests
- No mocking framework — tests construct `ast.FuncDecl` structs directly
- **Coverage**: 93.1% total (88.2% for `FindTestFuncs`, 71.4% for `isValidTestFunc`)
- `apply/` package has 0% coverage (just an `init()` call)

**Error handling**: Errors logged with `log.Printf()` and return nil/empty — no error propagation to callers. `ParseLabelExp` errors print to stdout and silently disable label filtering (graceful degradation).

## 3. Constraints (What limits decisions)

**Technical**:
- Go 1.26 required (module declares `go 1.26`)
- Package filtering is per-process — each `go test` package binary must independently equip gotest-labels
- Test name regex is the only mutation vector — same-named tests in different packages cannot be distinguished
- `defaultPkg` is a package-level `var` set to `"./..."` — tests override it (mutable global state)

**Quality**:
- CI runs `golangci-lint` v2.12.2 (1m timeout) + `go test -v -count=1 ./... -coverprofile`
- No `.golangci.yml` config file — uses defaults
- Codecov integration via CI secret
- No benchmarks found

**Operational**:
- CI: `go-ci.yaml` on PR/push to main (lint + test + codecov)
- Release: `go-pkg.yaml` on version tags (indexes on pkg.go.dev)
- No release automation beyond pkg.go.dev indexing

**Known limitation** (from README): Duplicate test function names across packages are both selected or both skipped. Mitigation: run packages separately.

## 4. Reusability (What to leverage)

**Similar implementations within codebase**:
- Three usage patterns are documented and exemplified:
  1. Anonymous import: `_ "github.com/maxwu/gotest-labels/apply"` (`examples/simple/`)
  2. Explicit `init()` call: `gotest_labels.MutateTestFilterByLabels()` (`examples/explicit/`)
  3. `TestMain()` call: safest for parent packages referencing sub-packages (`examples/explicittestmain/`)

**Key exported API surface**:
- `MutateTestFilterByLabels() map[string]TestLabels` — the only public entrypoint
- `TestLabels` = `map[string]string` — public type for label maps
- `ParseLabelExp(string) (Node, error)` — expression parser
- `Evaluate(Node, TestLabels) bool` — AST evaluator
- `FindTestFuncs([]string, Node) (map[string]TestLabels, error)` — AST-based test discovery

## 5. Handoff (What's next)

**For PLAN**: Key constraints when adding features:
- Any new flag must be both parsed in `cli.go` and stripped from `os.Args` (see `removeLabelFlags` pattern)
- The `defaultPkg` global var pattern may need refactoring for better testability
- `isValidTestFunc` only checks `*testing.T` — doesn't support `testing.B`, `testing.F` (fuzz), or subtests
- Expression parser could be extended for new operators (e.g., `!=`, `=~` regex match)
- Label values with spaces in comments are tokenized incorrectly (tokenizer splits on space)

**For CODE**:
- Test runner: `go test -v -count=1 ./...`
- Linter: `golangci-lint run --timeout 1m`
- Coverage: `go test -coverprofile=coverage.txt ./...`
- No `main.go` — this is a library, not a binary

**For COMMIT**:
- CI gates: lint passes + all tests pass + codecov upload
- PRs target `main` branch
- Tag format: `vX.Y.Z` triggers pkg.go.dev indexing
- License: Apache 2.0

**Gaps**:
- No `CLAUDE.md` with project-specific instructions
- `apply/` package has 0% test coverage
- `isValidTestFunc` at 71.4% — untested edge cases (non-star params, non-selector types)
- No fuzz tests or benchmarks
- `.history/` directory present but likely not git-tracked (local editor artifacts)
