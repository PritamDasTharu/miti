# miti

[![Go Reference](https://pkg.go.dev/badge/github.com/PritamDasTharu/miti.svg)](https://pkg.go.dev/github.com/PritamDasTharu/miti)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An ultra-lightweight, zero-dependency Bikram Sambat (BS / Nepali) ↔ Gregorian (AD) calendar library for Go.

- **Zero dependencies**: Built strictly with the Go standard library.
- **Zero allocations**: Core conversion (`FromTime`, `FromTimeNepal`, `ToTime`) and formatting paths allocate 0 heap bytes.
- **Micro footprint**: Static calendar lookup tables compile into $\approx 2\text{ KB}$ of static memory.
- **Data provenance**: Verified against official government gazettes, printed *Panchanga*, and independent golden anchors.
- **SQL & JSON safe**: Symmetric database driver support with SQL `NULL` and JSON `null` handling.

---

## Supported Range & Provenance

`miti` supports the continuous range **1970–2100 BS** (~131 years, ~47,847 days):
- **1970–2083 BS**: Officially published and verified against historical government records and *Nepal Panchanga Nirnayak Bikas Samiti* calendars.
- **2084–2100 BS**: Computed from astronomical calendar projections.

The code exposes `miti.MinYear (1970)`, `miti.MaxYear (2100)`, and `miti.MaxVerifiedYear` programmatically so critical financial or legal systems can validate bounds at runtime.

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
	// Current Nepali date (using zero-alloc Nepal timezone math)
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

## Contributing & Reporting Discrepancies

Found an inaccurate date? Please open an issue using our [Wrong Date Report](.github/ISSUE_TEMPLATE/wrong_date.md) template and cite official provenance.
