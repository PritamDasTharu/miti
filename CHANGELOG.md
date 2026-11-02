# Changelog

All notable changes to `miti` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-11-02
### Added
- **Core Conversion Engine**: $O(1)$ zero-allocation conversions between Bikram Sambat (BS) and Gregorian (AD) calendars.
- **Coverage & Provenance**: Continuous range covering 1970–2100 BS (~131 years, 47,849 days), with 1970–2083 BS officially verified against government records.
- **Micro-Footprint Table**: Static month lookup and cumulative offset tables compiling into $\approx 2.1\text{ KB}$ of static memory.
- **Data Verification**: Golden test suite checking 1 Baisakh across all 131 years, mid-year gazette anchors, and exhaustive roundtrip checking across all 47,849 days.
- **Timezone Safety**: Zero-alloc `FromTimeNepal` and `Today()` preventing UTC server day-shift traps.
- **Date Utilities & Arithmetic**: `StartOfMonth`, `EndOfMonth`, `IsWeekend`, `Compare`, `Sub`, `AddDays`, and `AddMonths` with end-of-month clamping.
- **Formatting & Parsing**: Token-based formatting (`Format`, `FormatIn`, `AppendFormat`, `AppendFormatIn`), Devanagari localization, and `[...]` literal escaping.
- **Fiscal Year**: Full support for Nepali fiscal years (Shrawan 1 start, Baisakh–Ashadh rules, century rollover `2099/00`, and quarters).
- **SQL & JSON Integrations**: Safe symmetric `driver.Valuer`/`Scanner` returning AD `time.Time`, plus opt-in `BSSQLDate` for BS text columns.
- **Examples & CI**: Isolated GORM and Fiber integration examples, zero-allocation assertions, and multi-OS GitHub Actions workflow.
