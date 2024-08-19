package common

func LowHigh(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	}
	return a, b
}
