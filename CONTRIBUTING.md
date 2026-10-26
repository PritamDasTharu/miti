# Contributing to miti

Thank you for your interest in contributing to `miti`!

## Calendar Accuracy & Bug Reports

The core priority of `miti` is **data trustworthiness**. If you spot an inaccurate date:
1. Open an issue using the [Wrong Date Report](.github/ISSUE_TEMPLATE/wrong_date.md) template.
2. Provide an authoritative citation (official *Nepal Panchanga Nirnayak Bikas Samiti* calendar, government gazette, or published national holiday notice).
3. Do not open PRs with arbitrary month length adjustments without citing official provenance.

## Development Guidelines

- **Zero Third-Party Dependencies**: The root module must strictly depend on the Go standard library only (`database/sql/driver`, `time`, etc.). No external modules in the root `go.mod`.
- **Zero Allocations**: Core conversion and formatting paths (`FromTime`, `FromTimeNepal`, `ToTime`, `AddDays`, `DaysInMonth`, `AppendFormat`) must maintain zero heap allocations.
- **Infallible Validated APIs**: Valid `Date` instances must never panic. Operations on valid dates that cannot fail must not return errors.
- **Code Generation**: Month tables and cumulative offsets in `table.go` are generated strictly from `data/bs_calendar.json` via `go generate ./...`. Never edit `table.go` manually.

## Running Tests

```bash
# Run standard unit tests & invariants
go test -v ./...

# Run race detector
go test -race ./...

# Verify code generation consistency
go generate ./... && git diff --exit-code
```
