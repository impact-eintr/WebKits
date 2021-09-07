package stdsort

func Search(n int, f func(int) bool) int {
	l, h := 0, n
	for l < h {
		mid := int(uint(l+h) >> 1)
		if !f(mid) {
			l = mid + 1
		} else {
			h = mid
		}
	}
	// l == h f(l-1) == false || f(h) == true
	return l
}
