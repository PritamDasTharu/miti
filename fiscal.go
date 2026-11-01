package miti

import (
	"strconv"
)

// FiscalYear returns the Nepali fiscal year (Aarthik Barsha) string (e.g. "2080/81").
//
// In Nepal, the fiscal year begins on 1 Shrawan (Month 4) and concludes on the final day
// of Ashadh (Month 3) of the following calendar year.
//
// Rules:
//   - Months 4..12 (Shrawan through Chaitra): fiscal year is "YYYY/(YY+1)" (e.g. 1 Shrawan 2083 -> "2083/84").
//   - Months 1..3 (Baisakh through Ashadh): belongs to the preceding fiscal year (e.g. 1 Baisakh 2083 -> "2082/83").
//   - Century rollover: year 2099 returns "2099/00".
//   - Returns "" for the zero Date{}.
func (d Date) FiscalYear() string {
	if !d.IsValid() {
		return ""
	}

	var startYear int
	var endYearTwoDigit int

	if d.month >= uint8(Shrawan) {
		startYear = int(d.year)
		endYearTwoDigit = (int(d.year) + 1) % 100
	} else {
		startYear = int(d.year) - 1
		endYearTwoDigit = int(d.year) % 100
	}

	var endStr string
	if endYearTwoDigit < 10 {
		endStr = "0" + strconv.Itoa(endYearTwoDigit)
	} else {
		endStr = strconv.Itoa(endYearTwoDigit)
	}

	return strconv.Itoa(startYear) + "/" + endStr
}

// FiscalYearRange returns the starting date (1 Shrawan) and ending date (last day of Ashadh)
// of the fiscal year containing d.
//
// Returns (Date{}, Date{}, ErrInvalidDate) if d is the zero Date{}.
// Returns (Date{}, Date{}, ErrOutOfRange) at table edges:
//   - Dates in Baisakh-Ashadh 1970 BS (whose fiscal year began in 1969 BS, before MinYear).
//   - Dates in Shrawan-Chaitra 2100 BS (whose fiscal year ends in 2101 BS, beyond MaxYear).
func (d Date) FiscalYearRange() (start Date, end Date, err error) {
	if !d.IsValid() {
		return Date{}, Date{}, ErrInvalidDate
	}

	var fyStartYear int
	var fyEndYear int

	if d.month >= uint8(Shrawan) {
		fyStartYear = int(d.year)
		fyEndYear = int(d.year) + 1
	} else {
		fyStartYear = int(d.year) - 1
		fyEndYear = int(d.year)
	}

	if fyStartYear < MinYear || fyEndYear > MaxYear {
		return Date{}, Date{}, ErrOutOfRange
	}

	// Start: 1 Shrawan fyStartYear
	startDate, err := New(fyStartYear, Shrawan, 1)
	if err != nil {
		return Date{}, Date{}, err
	}

	// End: last day of Ashadh fyEndYear
	ashadhDays := int(monthDays[fyEndYear-MinYear][Ashadh-1])
	endDate, err := New(fyEndYear, Ashadh, ashadhDays)
	if err != nil {
		return Date{}, Date{}, err
	}

	return startDate, endDate, nil
}

// Quarter returns the Nepali fiscal quarter (1..4) for d:
//   - Q1: Shrawan, Bhadra, Ashoj (Months 4..6)
//   - Q2: Kartik, Mangsir, Poush (Months 7..9)
//   - Q3: Magh, Falgun, Chaitra (Months 10..12)
//   - Q4: Baisakh, Jestha, Ashadh (Months 1..3)
//
// Returns 0 for the zero Date{}.
func (d Date) Quarter() int {
	if !d.IsValid() {
		return 0
	}
	switch d.month {
	case uint8(Shrawan), uint8(Bhadra), uint8(Ashoj):
		return 1
	case uint8(Kartik), uint8(Mangsir), uint8(Poush):
		return 2
	case uint8(Magh), uint8(Falgun), uint8(Chaitra):
		return 3
	default:
		return 4
	}
}
