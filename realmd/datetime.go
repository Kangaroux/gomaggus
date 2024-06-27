package realmd

import "time"

// DateTime stores a datetime relative to Jan 1, 2000. The datetime is represented in binary
// as a compressed datetime object. Bit structure:
//
//	0..5		minute
//	6..10		hour
//	11..13		day of the week (0 = sunday)
//	14..19		day of the month (0 = 1st of the month)
//	20..23		month (0 = january)
//	24..28		year (relative to 2000)
//	29..32		unused
type DateTime uint32

func NewDateTime(t time.Time) DateTime {
	return DateTime(t.Minute()) |
		DateTime(t.Hour())<<6 |
		DateTime(t.Weekday())<<11 |
		DateTime(t.Day()-1)<<14 |
		DateTime(t.Month()-1)<<20 |
		DateTime(t.Year()-2000)<<24
}
