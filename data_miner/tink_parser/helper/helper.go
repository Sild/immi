package helper

import "time"

func FormatDate(ts time.Time) string {
	return ts.Format("2006-01-02 -0700")
}

func FormatTime(ts time.Time) string {
	return ts.Format("2006-01-02 15:04:05 -0700")

}

func DateFromStr(str string) time.Time {
	t, _ := time.Parse("2006-01-02", str)
	return t
}

func DBDateFromStr(str string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", str)
	return t
}

func TimeFromStr(str string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", str)
	return t
}
