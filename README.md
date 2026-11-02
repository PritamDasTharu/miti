# miti

[![Go Reference](https://pkg.go.dev/badge/github.com/PritamDasTharu/miti.svg)](https://pkg.go.dev/github.com/PritamDasTharu/miti)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/PritamDasTharu/miti/actions/workflows/ci.yml/badge.svg)](https://github.com/PritamDasTharu/miti/actions/workflows/ci.yml)

An ultra-lightweight, zero-dependency Bikram Sambat (BS / Nepali) ↔ Gregorian (AD) calendar library for Go.

- **Zero dependencies**: Strictly Go standard library only (`time`, `database/sql/driver`, `encoding/json`).
- **Zero allocations**: Core conversion (`FromTime`, `FromTimeNepal`, `ToTime`), arithmetic (`AddDays`), and formatting (`AppendFormat`) allocate **0 heap bytes**.
- **Micro footprint**: Static lookup tables compile into $\approx 2.1\text{ KB}$ of static memory.
- **Verified provenance**: Cross-verified against official government gazettes, printed *Panchanga*, and independent golden anchors.
- **Database safety**: Symmetric SQL driver support with standard SQL `DATE` and `VARCHAR` options.

---

## Benchmarks

Measured on AMD Ryzen 7 5700U (Go 1.26 / amd64):

| Benchmark | Operations/sec | Speed (ns/op) | Memory (B/op) | Allocations/op |
|:---|---:|---:|---:|---:|
| `ToTime` (BS $\to$ AD) | **212,245,425** | **6.6 ns** | **0 B** | **0 allocs** |
| `AddDays` | **82,739,218** | **15.8 ns** | **0 B** | **0 allocs** |
| `FromTime` (AD $\to$ BS) | **29,395,480** | **45.2 ns** | **0 B** | **0 allocs** |
| `FromTimeNepal` (UTC $\to$ NPT) | **22,505,878** | **57.4 ns** | **0 B** | **0 allocs** |
| `AppendFormat` | **11,116,009** | **91.8 ns** | **0 B** | **0 allocs** |

---

## Supported Range & Provenance

`miti` covers the continuous range **1970–2100 BS** (~131 years, 47,849 days):
- **1970–2083 BS**: Officially published and verified against historical government records and *Nepal Panchanga Nirnayak Bikas Samiti* calendars.
- **2084–2100 BS**: Computed from astronomical calendar projection models.

The code exposes `miti.MinYear (1970)`, `miti.MaxYear (2100)`, and `miti.MaxVerifiedYear (2083)` programmatically so critical financial or legal systems can validate bounds at runtime.

---

## Quickstart

```go
package main

import (
	"fmt"
	"time"

	"github.com/PritamDasTharu/miti"
)

func main() {
	// Current Nepali date (zero-alloc Nepal timezone math)
	today := miti.Today()
	fmt.Println("Today in Nepal:", today) // e.g. 2083-07-18

	// Convert AD to BS
	adTime := time.Date(2015, time.September, 20, 0, 0, 0, 0, time.UTC)
	bsDate, err := miti.FromTime(adTime)
	if err != nil {
		panic(err)
	}
	fmt.Println("Constitution Day BS:", bsDate) // 2072-06-03

	// Convert BS to AD (returns midnight UTC)
	adResult := bsDate.ToTime()
	fmt.Println("Back to AD:", adResult.Format("2006-01-02")) // 2015-09-20
}
```

> [!WARNING]
> **UTC Server Day-Shift Trap**:
> Most cloud servers, Docker containers, and CI pipelines run with their system clock in **UTC**.
> Because Nepal Standard Time is `UTC+05:45`, between **18:15 UTC and 24:00 UTC**, Nepal is already on the **next calendar day**.
> Calling `miti.FromTime(time.Now())` on a UTC server will yield yesterday's Nepali date during those 5 hours and 45 minutes!
>
> To convert a server's wall-clock instant to the current Nepali date safely without allocating timezone objects, always use:
> ```go
> bsDate, err := miti.FromTimeNepal(time.Now())
> // or simply:
> today := miti.Today()
> ```

---

## Formatting and Parsing

`miti` provides intuitive token-based formatting with Devanagari localization and bracketed literal escaping `[...]`:

```go
d, _ := miti.New(2081, miti.Baisakh, 1)

// English layout
fmt.Println(d.Format("dddd, D MMMM YYYY"))
// Output: Saturday, 1 Baisakh 2081

// Devanagari numerals and localized names
fmt.Println(d.FormatIn("dddd, D MMMM YYYY", miti.LocaleNepali))
// Output: शनिबार, १ वैशाख २०८१

// Literal text escaping using [...]
fmt.Println(d.Format("[Date: ]YYYY-MM-DD"))
// Output: Date: 2081-01-01

// Parsing English and Devanagari dates
parsed, err := miti.Parse("YYYY-MM-DD", "२०८१-०१-०१")
```

---

## Nepali Fiscal Year (Aarthik Barsha)

Nepali fiscal years start on **1 Shrawan** (Month 4) and end on the last day of **Ashadh** (Month 3) of the subsequent year:

```go
d, _ := miti.New(2081, miti.Baisakh, 1)

fmt.Println(d.FiscalYear()) // "2080/81" (Baisakh belongs to preceding FY)
fmt.Println(d.Quarter())    // 4 (Q4: Baisakh-Ashadh)

start, end, err := d.FiscalYearRange()
// start = 2080-04-01, end = 2081-03-31
```

---

## Database & JSON Integrations

### 1. SQL `DATE` Columns (Default)
By default, `miti.Date` returns an AD `time.Time` at midnight UTC from `Value()`, and `Scan()` parses AD representations. This avoids corrupting standard SQL `DATE` columns:

```go
type Invoice struct {
    ID   uint       `gorm:"primaryKey"`
    Date miti.Date  `gorm:"type:date;not null"`
}
```

> [!NOTE]
> In PostgreSQL, always use `DATE`, not `TIMESTAMPTZ`.

### 2. SQL `VARCHAR` / `TEXT` Columns (Opt-in)
If you specifically want to store Bikram Sambat text strings (`"2081-01-15"`):

```go
type Document struct {
    ID     uint            `gorm:"primaryKey"`
    DateBS miti.BSSQLDate  `gorm:"type:varchar(10);not null"`
}
```

---

## Contributing & Reporting Discrepancies

Found an inaccurate date? Please open an issue using our [Wrong Date Report](.github/ISSUE_TEMPLATE/wrong_date.md) template and cite official provenance.
