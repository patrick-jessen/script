package pe

import "math"

func multipleOf(val int32, of int32) (r int32) {
	r = val
	if val%of > 0 {
		r = (val/of + 1) * of
	}
	return
}

func isPow2(num int) bool {
	ln := math.Log2(float64(num))
	return math.Ceil(ln) == math.Floor(ln)
}
