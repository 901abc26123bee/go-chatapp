package timeutil

import "time"

// defines default layout of time
const (
	DefaultTimeLayout = time.RFC3339

	DateLayout     = "20060102"
	DateTimeLayout = "2006-01-02 15:04:05"
)

// TimezoneTaipei defines timezone of taipei
var TimezoneTaipei = time.FixedZone("Asia/Taipei", 8*3600)

// ConvertUTCTimeISOString convert the time with UTC time zone into iso 8601 format string
func ConvertUTCTimeISOString(t time.Time) string {
	return t.UTC().Format(DefaultTimeLayout)
}

// ConvertTimeToString convert the time with specified time zone into specified layout string
func ConvertTimeToString(t time.Time, timezone *time.Location, layout string) string {
	return t.In(timezone).Format(layout)
}
