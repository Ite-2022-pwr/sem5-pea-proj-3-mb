package utils

func SecondsToNanoSeconds(seconds int64) int64 {
	return seconds * 1e9
}

func MinutesToNanoSeconds(minutes int64) int64 {
	return SecondsToNanoSeconds(minutes * 60)
}
