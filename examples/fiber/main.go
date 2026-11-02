package main

import (
	"encoding/json"
	"fmt"

	"github.com/PritamDasTharu/miti"
)

// SchedulePayload demonstrates Fiber request/response JSON payload.
type SchedulePayload struct {
	Title     string    `json:"title"`
	StartDate miti.Date `json:"start_date"`
	EndDate   miti.Date `json:"end_date"`
}

func main() {
	start, _ := miti.New(2081, miti.Baisakh, 1)
	end, _ := miti.New(2081, miti.Baisakh, 31)

	payload := SchedulePayload{
		Title:     "HydroCMMS Maintenance Sprint",
		StartDate: start,
		EndDate:   end,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("Serialized JSON Payload:")
	fmt.Println(string(data))

	var parsed SchedulePayload
	if err := json.Unmarshal(data, &parsed); err != nil {
		panic(err)
	}
	fmt.Printf("Parsed Dates: Start=%v, End=%v, DaysBetween=%d\n",
		parsed.StartDate, parsed.EndDate, parsed.EndDate.Sub(parsed.StartDate))
}
