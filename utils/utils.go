package utils

import log "github.com/sirupsen/logrus"

func PanicIf(err error) {
	if err != nil {
		log.Panic()
	}
}

func member(n int, set []int) bool {
	for m := range set {
		if n == m {
			return true
		}
	}
	return false
}

func Subset(subset []int, superset []int) bool {
	for n := range subset {
		if !member(n, superset) {
			return false
		}
	}
	return true
}

func Merge(set1 []int, set2 []int) []int {
	ret := make([]int, len(set1))
	copy(ret, set1)
	for n := range set2 {
		if !member(n, ret) {
			ret = append(ret, n)
		}
	}
	return ret
}