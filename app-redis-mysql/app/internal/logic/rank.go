package logic

import "strings"

func Rank(searchTerm string, words []string, max int) *[]string {
	capacity := Min(max, len(words))
	top := make([]string, 0, capacity)
	leftovers := []string{}

	for _, w := range words {
		if strings.HasPrefix(w, searchTerm) {
			top = append(top, w)
		} else {
			leftovers = append(leftovers, w)
		}
		if len(top) == capacity {
			return &top
		}
	}
	words = leftovers
	leftovers = []string{}
	for _, w := range words {
		if strings.Contains(w, searchTerm) {
			top = append(top, w)
		} else {
			leftovers = append(leftovers, w)
		}
		if len(top) == capacity {
			return &top
		}
	}
	leftovers = leftovers[0:(capacity - len(top))]
	top = append(top, leftovers...)
	return &top
}

func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
