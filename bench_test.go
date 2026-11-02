package miti_test

import (
	"testing"
	"time"

	"github.com/PritamDasTharu/miti"
)

func BenchmarkFromTime(b *testing.B) {
	now := time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = miti.FromTime(now)
	}
}

func BenchmarkFromTimeNepal(b *testing.B) {
	now := time.Date(2024, 4, 27, 18, 30, 0, 0, time.UTC)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = miti.FromTimeNepal(now)
	}
}

func BenchmarkToTime(b *testing.B) {
	d, _ := miti.New(2081, miti.Baisakh, 15)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.ToTime()
	}
}

func BenchmarkAddDays(b *testing.B) {
	d, _ := miti.New(2081, miti.Baisakh, 15)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = d.AddDays(30)
	}
}

func BenchmarkAppendFormat(b *testing.B) {
	d, _ := miti.New(2081, miti.Baisakh, 15)
	buf := make([]byte, 0, 32)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bSlice := buf[:0]
		_ = d.AppendFormat(bSlice, "YYYY-MM-DD")
	}
}
