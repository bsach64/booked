package utils

import "time"

func GetUTCTimeFromUnix(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}

func GetUTCTime(t time.Time) time.Time {
	return t.UTC()
}

func GetUTCUnixTime(t time.Time) int64 {
	return t.UTC().Unix()
}
