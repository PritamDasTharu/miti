package miti_test

import (
	"fmt"
	"time"

	"github.com/PritamDasTharu/miti"
)

func ExampleFromTime() {
	adTime := time.Date(2015, time.September, 20, 0, 0, 0, 0, time.UTC)
	bsDate, err := miti.FromTime(adTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(bsDate)
	// Output:
	// 2072-06-03
}

func ExampleDate_ToTime() {
	d, _ := miti.New(2072, miti.Ashoj, 3)
	adTime := d.ToTime()
	fmt.Println(adTime.Format("2006-01-02"))
	// Output:
	// 2015-09-20
}

func ExampleDate_FiscalYear() {
	d1, _ := miti.New(2080, miti.Shrawan, 1)
	d2, _ := miti.New(2081, miti.Baisakh, 1)

	fmt.Println(d1.FiscalYear())
	fmt.Println(d2.FiscalYear())
	// Output:
	// 2080/81
	// 2080/81
}

func ExampleDate_Format() {
	d, _ := miti.New(2081, miti.Baisakh, 1)

	fmt.Println(d.Format("dddd, D MMMM YYYY"))
	fmt.Println(d.FormatIn("dddd, D MMMM YYYY", miti.LocaleNepali))
	// Output:
	// Saturday, 1 Baisakh 2081
	// शनिबार, १ वैशाख २०८१
}

func ExampleParse() {
	d, err := miti.Parse("YYYY-MM-DD", "2081-01-01")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d-%02d-%02d\n", d.Year(), d.Month(), d.Day())
	// Output:
	// 2081-01-01
}
