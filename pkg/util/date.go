package util

import "time"

const dateTimeFormat = "02.01.2006 15:04"
const dateFormat = "01.02.2006"

func FormatDate(date string, withTime bool) string {
	if date == "" {
		return ""
	}

	val, _ := time.Parse(time.RFC3339, date)
	if withTime {
		return val.Format(dateTimeFormat)
	}

	return val.Format(dateFormat)
}
