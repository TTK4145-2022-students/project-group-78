package utils

func Member(n int, set []int) bool {
	for m := range set {
		if n == m {
			return true
		}
	}
	return false
}

func Subset(subset []int, superset []int) bool {
	for n := range subset {
		if !Member(n, superset) {
			return false
		}
	}
	return true
}
