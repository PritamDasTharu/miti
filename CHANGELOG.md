# Changelog

All notable changes to `miti` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Core Bikram Sambat date conversion engine supporting 1970–2100 BS.
- Static lookup table code generation via `internal/gen/gen_table.go`.
- Golden tests validating 1 Baisakh across 131 years and mid-year historical anchors.
- Exhaustive day-by-day roundtrip suite covering all 47,849 days.
- Zero-alloc Nepal time conversion `FromTimeNepal` and `Today()`.
- Repository initialization and Go module definition (`go 1.21`).
