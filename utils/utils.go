package utils

func Member(n byte, set []byte) bool {
	for _, m := range set {
		if n == m {
			return true
		}
	}
	return false
}

func Subset(subset []byte, superset []byte) bool {
	for _, n := range subset {
		if !Member(n, superset) {
			return false
		}
	}
	return true
}

func Equal(s1 []byte, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
