package main

import (
	"fmt"

	"github.com/PritamDasTharu/miti"
)

// Invoice represents a business invoice model.
// Date is stored as standard SQL DATE (AD time.Time under the hood via driver.Valuer).
// FiscalYearBS uses BSSQLDate if raw BS string representation is desired in VARCHAR.
type Invoice struct {
	ID        uint           `gorm:"primaryKey"`
	Number    string         `gorm:"type:varchar(50);not null"`
	Date      miti.Date      `gorm:"type:date;not null"`
	BSString  miti.BSSQLDate `gorm:"type:varchar(10)"`
	TotalAmt  int64          `gorm:"not null"` // Amount in paisa
}

func main() {
	today := miti.Today()
	inv := Invoice{
		ID:       101,
		Number:   "INV-2083-001",
		Date:     today,
		BSString: miti.BSSQLDate(today),
		TotalAmt: 1250000, // 12,500.00 Rs
	}

	fmt.Printf("Created Invoice: %+v\n", inv)
	fmt.Printf("AD SQL Value (DATE): %v\n", inv.Date.ToTime().Format("2006-01-02"))
	fmt.Printf("BS Text Value (VARCHAR): %v\n", inv.BSString.Date().String())
	fmt.Printf("Fiscal Year: %s (Q%d)\n", inv.Date.FiscalYear(), inv.Date.Quarter())
}
