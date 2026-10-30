package miti

import (
	"strconv"
	"strings"
	"time"
)

// Locale specifies the language and numeral system used for formatting and parsing.
type Locale uint8

const (
	// LocaleEnglish formats with Latin numerals and English month/weekday names.
	LocaleEnglish Locale = iota
	// LocaleNepali formats with Devanagari numerals (०-९) and Nepali month/weekday names.
	LocaleNepali
)

// Standard layout presets.
const (
	StandardDate = "YYYY-MM-DD"
)

var devanagariDigits = [10]string{"०", "१", "२", "३", "४", "५", "६", "७", "८", "९"}

var englishMonths = [12]string{
	"Baisakh", "Jestha", "Ashadh", "Shrawan", "Bhadra", "Ashoj",
	"Kartik", "Mangsir", "Poush", "Magh", "Falgun", "Chaitra",
}

var englishShortMonths = [12]string{
	"Bai", "Jes", "Ash", "Shr", "Bha", "Ash",
	"Kar", "Man", "Pou", "Mag", "Fal", "Cha",
}

var nepaliMonths = [12]string{
	"वैशाख", "ज्येष्ठ", "आषाढ", "श्रावण", "भाद्र", "आश्विन",
	"कार्तिक", "मंसिर", "पौष", "माघ", "फाल्गुन", "चैत्र",
}

var nepaliShortMonths = [12]string{
	"वै", "ज्ये", "आषा", "श्रा", "भा", "आश्वि",
	"का", "मं", "पौ", "मा", "फा", "चै",
}

var englishWeekdays = [7]string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

var englishShortWeekdays = [7]string{
	"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat",
}

var nepaliWeekdays = [7]string{
	"आइतबार", "सोमबार", "मंगलबार", "बुधबार", "बिहिबार", "शुक्रबार", "शनिबार",
}

var nepaliShortWeekdays = [7]string{
	"आइत", "सोम", "मंगल", "बुध", "बिहि", "शुक्र", "शनि",
}

// monthAliases maps transliterations and Devanagari spellings to Month (1..12).
var monthAliases = map[string]Month{
	// English canonical & common romanizations
	"baisakh": Baisakh, "baishakh": Baisakh, "baisak": Baisakh, "bai": Baisakh,
	"jestha": Jestha, "jyestha": Jestha, "jeth": Jestha, "jes": Jestha,
	"ashadh": Ashadh, "asad": Ashadh, "ashad": Ashadh, "asar": Ashadh, "ash": Ashadh,
	"shrawan": Shrawan, "shravan": Shrawan, "saun": Shrawan, "shr": Shrawan,
	"bhadra": Bhadra, "bhadrapad": Bhadra, "bhadau": Bhadra, "bhadrupada": Bhadra, "bha": Bhadra,
	"ashoj": Ashoj, "asoj": Ashoj, "ashwin": Ashoj, "aswin": Ashoj,
	"kartik": Kartik, "katik": Kartik, "kar": Kartik,
	"mangsir": Mangsir, "mamsir": Mangsir, "margashirsha": Mangsir, "marga": Mangsir, "man": Mangsir,
	"poush": Poush, "paush": Poush, "push": Poush, "pous": Poush, "pou": Poush,
	"magh": Magh, "mag": Magh,
	"falgun": Falgun, "phagun": Falgun, "fagun": Falgun, "phalgun": Falgun, "fal": Falgun,
	"chaitra": Chaitra, "chait": Chaitra, "cha": Chaitra,

	// Devanagari canonical & spoken variants
	"वैशाख": Baisakh, "बैशाख": Baisakh, "वै": Baisakh,
	"ज्येष्ठ": Jestha, "जेठ": Jestha, "ज्ये": Jestha,
	"आषाढ": Ashadh, "असार": Ashadh, "आषा": Ashadh,
	"श्रावण": Shrawan, "साउन": Shrawan, "श्रा": Shrawan,
	"भाद्र": Bhadra, "भदौ": Bhadra, "भा": Bhadra,
	"आश्विन": Ashoj, "असोज": Ashoj, "आश्वि": Ashoj,
	"कार्तिक": Kartik, "कात्तिक": Kartik, "का": Kartik,
	"मार्गशीर्ष": Mangsir, "मंसिर": Mangsir, "मार्ग": Mangsir, "मं": Mangsir,
	"पौष": Poush, "पुस": Poush, "पौ": Poush,
	"माघ": Magh, "मा": Magh,
	"फाल्गुन": Falgun, "फागुन": Falgun, "फा": Falgun,
	"चैत्र": Chaitra, "चैत": Chaitra, "चै": Chaitra,
}

var weekdayAliases = map[string]time.Weekday{
	"sunday": time.Sunday, "sun": time.Sunday,
	"monday": time.Monday, "mon": time.Monday,
	"tuesday": time.Tuesday, "tue": time.Tuesday,
	"wednesday": time.Wednesday, "wed": time.Wednesday,
	"thursday": time.Thursday, "thu": time.Thursday,
	"friday": time.Friday, "fri": time.Friday,
	"saturday": time.Saturday, "sat": time.Saturday,

	"आइतबार": time.Sunday, "आइत": time.Sunday,
	"सोमबार": time.Monday, "सोम": time.Monday,
	"मंगलबार": time.Tuesday, "मंगल": time.Tuesday,
	"बुधबार": time.Wednesday, "बुध": time.Wednesday,
	"बिहिबार": time.Thursday, "बिहि": time.Thursday,
	"शुक्रबार": time.Friday, "शुक्र": time.Friday,
	"शनिबार": time.Saturday, "शनि": time.Saturday,
}

// appendDevanagariInt appends the Devanagari representation of n (with minDigits zero-padding) to dst.
func appendDevanagariInt(dst []byte, n int, minDigits int) []byte {
	var s string
	if minDigits == 2 && n < 10 {
		s = "0" + strconv.Itoa(n)
	} else if minDigits == 4 {
		s = fmtSprintf4Digit(n)
	} else {
		s = strconv.Itoa(n)
	}

	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			dst = append(dst, devanagariDigits[ch-'0']...)
		} else {
			dst = append(dst, byte(ch))
		}
	}
	return dst
}

func fmtSprintf4Digit(n int) string {
	var b [4]byte
	b[0] = byte('0' + (n/1000)%10)
	b[1] = byte('0' + (n/100)%10)
	b[2] = byte('0' + (n/10)%10)
	b[3] = byte('0' + n%10)
	return string(b[:])
}

// Format formats d according to the layout string using LocaleEnglish.
// Returns "" for the zero Date{}.
func (d Date) Format(layout string) string {
	return d.FormatIn(layout, LocaleEnglish)
}

// FormatIn formats d according to the layout string using the specified locale.
// Returns "" for the zero Date{}.
func (d Date) FormatIn(layout string, loc Locale) string {
	if !d.IsValid() {
		return ""
	}
	// Estimate capacity
	buf := make([]byte, 0, len(layout)+16)
	buf = d.AppendFormatIn(buf, layout, loc)
	return string(buf)
}

// AppendFormat appends the formatted date to dst using LocaleEnglish.
// Returns dst unchanged for the zero Date{}.
func (d Date) AppendFormat(dst []byte, layout string) []byte {
	return d.AppendFormatIn(dst, layout, LocaleEnglish)
}

// AppendFormatIn appends the formatted date to dst using the specified locale.
// Returns dst unchanged for the zero Date{}.
func (d Date) AppendFormatIn(dst []byte, layout string, loc Locale) []byte {
	if !d.IsValid() {
		return dst
	}

	inEscape := false
	i := 0
	for i < len(layout) {
		ch := layout[i]

		// Bracket escape syntax [...]
		if ch == '[' {
			inEscape = true
			i++
			continue
		}
		if ch == ']' && inEscape {
			inEscape = false
			i++
			continue
		}
		if inEscape {
			dst = append(dst, ch)
			i++
			continue
		}

		// Token checking
		rem := layout[i:]
		switch {
		case strings.HasPrefix(rem, "YYYY"):
			y := int(d.year)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, y, 4)
			} else {
				dst = append(dst,
					byte('0'+(y/1000)%10),
					byte('0'+(y/100)%10),
					byte('0'+(y/10)%10),
					byte('0'+y%10),
				)
			}
			i += 4

		case strings.HasPrefix(rem, "YY"):
			// Format-only 2-digit year
			y2 := int(d.year % 100)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, y2, 2)
			} else {
				dst = append(dst,
					byte('0'+(y2/10)%10),
					byte('0'+y2%10),
				)
			}
			i += 2

		case strings.HasPrefix(rem, "MMMM"):
			mIdx := int(d.month) - 1
			if loc == LocaleNepali {
				dst = append(dst, nepaliMonths[mIdx]...)
			} else {
				dst = append(dst, englishMonths[mIdx]...)
			}
			i += 4

		case strings.HasPrefix(rem, "MMM"):
			mIdx := int(d.month) - 1
			if loc == LocaleNepali {
				dst = append(dst, nepaliShortMonths[mIdx]...)
			} else {
				dst = append(dst, englishShortMonths[mIdx]...)
			}
			i += 3

		case strings.HasPrefix(rem, "MM"):
			m := int(d.month)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, m, 2)
			} else {
				dst = append(dst, byte('0'+m/10), byte('0'+m%10))
			}
			i += 2

		case strings.HasPrefix(rem, "M"):
			m := int(d.month)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, m, 1)
			} else {
				if m >= 10 {
					dst = append(dst, '1', byte('0'+m%10))
				} else {
					dst = append(dst, byte('0'+m))
				}
			}
			i += 1

		case strings.HasPrefix(rem, "dddd"):
			wIdx := int(d.Weekday())
			if loc == LocaleNepali {
				dst = append(dst, nepaliWeekdays[wIdx]...)
			} else {
				dst = append(dst, englishWeekdays[wIdx]...)
			}
			i += 4

		case strings.HasPrefix(rem, "ddd"):
			wIdx := int(d.Weekday())
			if loc == LocaleNepali {
				dst = append(dst, nepaliShortWeekdays[wIdx]...)
			} else {
				dst = append(dst, englishShortWeekdays[wIdx]...)
			}
			i += 3

		case strings.HasPrefix(rem, "DD"):
			day := int(d.day)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, day, 2)
			} else {
				dst = append(dst, byte('0'+day/10), byte('0'+day%10))
			}
			i += 2

		case strings.HasPrefix(rem, "D"):
			day := int(d.day)
			if loc == LocaleNepali {
				dst = appendDevanagariInt(dst, day, 1)
			} else {
				if day >= 10 {
					dst = append(dst, byte('0'+day/10), byte('0'+day%10))
				} else {
					dst = append(dst, byte('0'+day))
				}
			}
			i += 1

		default:
			dst = append(dst, ch)
			i++
		}
	}

	return dst
}

// Parse parses a formatted string value into a Date according to layout.
// It supports both English and Devanagari numerals and month aliases.
//
// Rules:
//   - "YY" is format-only and returns ErrInvalidFormat if present in layout.
//   - Unknown token runs (e.g. MMMMM, YYY, dd, lone d) return ErrInvalidFormat.
//   - Adjacent variable-width tokens without delimiters (e.g. MD, DM) return ErrInvalidFormat.
//   - Weekday tokens (dddd, ddd) are actively validated against the date; mismatch returns ErrInvalidFormat.
func Parse(layout, value string) (Date, error) {
	if len(layout) == 0 || len(value) == 0 {
		return Date{}, ErrInvalidFormat
	}

	// 1. Reject YY format-only token and unknown token patterns in layout
	if err := validateLayoutForParsing(layout); err != nil {
		return Date{}, err
	}

	var year int = -1
	var month Month = 0
	var day int = -1
	var parsedWeekday time.Weekday = time.Weekday(99)

	lIdx := 0
	vIdx := 0
	vRunes := []rune(value)

	for lIdx < len(layout) {
		ch := layout[lIdx]

		// Handle bracket escape [...] in layout
		if ch == '[' {
			lIdx++
			for lIdx < len(layout) && layout[lIdx] != ']' {
				if vIdx >= len(vRunes) || vRunes[vIdx] != rune(layout[lIdx]) {
					return Date{}, ErrInvalidFormat
				}
				vIdx++
				lIdx++
			}
			if lIdx < len(layout) && layout[lIdx] == ']' {
				lIdx++
			}
			continue
		}

		remLayout := layout[lIdx:]
		switch {
		case strings.HasPrefix(remLayout, "YYYY"):
			lIdx += 4
			digits, advance := readDigits(vRunes[vIdx:], 4, 4)
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			year = digits
			vIdx += advance

		case strings.HasPrefix(remLayout, "MMMM") || strings.HasPrefix(remLayout, "MMM"):
			isFull := strings.HasPrefix(remLayout, "MMMM")
			if isFull {
				lIdx += 4
			} else {
				lIdx += 3
			}
			m, advance := readMonthName(vRunes[vIdx:])
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			month = m
			vIdx += advance

		case strings.HasPrefix(remLayout, "MM"):
			lIdx += 2
			digits, advance := readDigits(vRunes[vIdx:], 2, 2)
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			month = Month(digits)
			vIdx += advance

		case strings.HasPrefix(remLayout, "M"):
			lIdx += 1
			digits, advance := readDigits(vRunes[vIdx:], 1, 2)
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			month = Month(digits)
			vIdx += advance

		case strings.HasPrefix(remLayout, "dddd") || strings.HasPrefix(remLayout, "ddd"):
			isFull := strings.HasPrefix(remLayout, "dddd")
			if isFull {
				lIdx += 4
			} else {
				lIdx += 3
			}
			w, advance := readWeekdayName(vRunes[vIdx:])
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			parsedWeekday = w
			vIdx += advance

		case strings.HasPrefix(remLayout, "DD"):
			lIdx += 2
			digits, advance := readDigits(vRunes[vIdx:], 2, 2)
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			day = digits
			vIdx += advance

		case strings.HasPrefix(remLayout, "D"):
			lIdx += 1
			digits, advance := readDigits(vRunes[vIdx:], 1, 2)
			if advance == 0 {
				return Date{}, ErrInvalidFormat
			}
			day = digits
			vIdx += advance

		default:
			if vIdx >= len(vRunes) || vRunes[vIdx] != rune(ch) {
				return Date{}, ErrInvalidFormat
			}
			lIdx++
			vIdx++
		}
	}

	if vIdx != len(vRunes) {
		return Date{}, ErrInvalidFormat
	}

	if year == -1 || month == 0 || day == -1 {
		return Date{}, ErrInvalidFormat
	}

	d, err := New(year, month, day)
	if err != nil {
		return Date{}, err
	}

	// Validate weekday if present in layout
	if parsedWeekday != time.Weekday(99) {
		if d.Weekday() != parsedWeekday {
			return Date{}, ErrInvalidFormat
		}
	}

	return d, nil
}

func validateLayoutForParsing(layout string) error {
	inEscape := false
	i := 0
	for i < len(layout) {
		ch := layout[i]
		if ch == '[' {
			inEscape = true
			i++
			continue
		}
		if ch == ']' {
			inEscape = false
			i++
			continue
		}
		if inEscape {
			i++
			continue
		}

		rem := layout[i:]

		// Reject YY in parsing (ambiguous across two centuries)
		if strings.HasPrefix(rem, "YY") && !strings.HasPrefix(rem, "YYYY") {
			return ErrInvalidFormat
		}

		// Reject adjacent variable-width tokens like MD or DM without delimiter
		if (strings.HasPrefix(rem, "MD") && !strings.HasPrefix(rem, "MMM")) ||
			strings.HasPrefix(rem, "DM") {
			return ErrInvalidFormat
		}

		// Reject unknown runs like MMMMM, YYY, dd, lone d
		if strings.HasPrefix(rem, "MMMMM") ||
			(strings.HasPrefix(rem, "YYY") && !strings.HasPrefix(rem, "YYYY")) ||
			(strings.HasPrefix(rem, "dd") && !strings.HasPrefix(rem, "ddd")) {
			return ErrInvalidFormat
		}

		// Step forward
		switch {
		case strings.HasPrefix(rem, "YYYY"):
			i += 4
		case strings.HasPrefix(rem, "MMMM"):
			i += 4
		case strings.HasPrefix(rem, "MMM"):
			i += 3
		case strings.HasPrefix(rem, "MM"):
			i += 2
		case strings.HasPrefix(rem, "M"):
			i += 1
		case strings.HasPrefix(rem, "dddd"):
			i += 4
		case strings.HasPrefix(rem, "ddd"):
			i += 3
		case strings.HasPrefix(rem, "DD"):
			i += 2
		case strings.HasPrefix(rem, "D"):
			i += 1
		default:
			i++
		}
	}
	return nil
}

func readDigits(runes []rune, minDigits, maxDigits int) (val int, advance int) {
	count := 0
	for count < maxDigits && count < len(runes) {
		r := runes[count]
		digit := -1
		if r >= '0' && r <= '9' {
			digit = int(r - '0')
		} else if r >= '०' && r <= '९' {
			digit = int(r - '०')
		} else {
			break
		}
		val = val*10 + digit
		count++
	}

	if count < minDigits {
		return 0, 0
	}
	return val, count
}

func readMonthName(runes []rune) (Month, int) {
	// Greedy prefix matching against monthAliases
	for length := min(20, len(runes)); length >= 1; length-- {
		candidate := strings.ToLower(string(runes[:length]))
		if m, ok := monthAliases[candidate]; ok {
			return m, length
		}
	}
	return 0, 0
}

func readWeekdayName(runes []rune) (time.Weekday, int) {
	for length := min(20, len(runes)); length >= 1; length-- {
		candidate := strings.ToLower(string(runes[:length]))
		if w, ok := weekdayAliases[candidate]; ok {
			return w, length
		}
	}
	return time.Sunday, 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

