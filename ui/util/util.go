package util

func Max(first int, second int) int {
	if first >= second {
		return first
	}
	return second
}

func Min(first int, second int) int {
	if first <= second {
		return first
	}
	return second
}
