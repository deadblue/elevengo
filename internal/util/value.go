package util

func NonZero(number ...int) int {
	for _, n := range number {
		if n != 0 {
			return n
		}
	}
	return 0
}

func NonEmptyString(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}
