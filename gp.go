package lufact

// lusolv solves a square linear system, given an LU factorization.
//
// Solve for X in a square linear system Ax = b, given the factorization
// PA=LU.
//
// Input parameters:
//   n                          dimension of matrix.
//   lu, lurow, lcolst,
//   ucolst, perm               PA=LU factorization (see lufact for format).
//
// Modified parameter:
//   x                          Real array of length n.
//                              On entry, holds B.  On exit, holds X.
//
// Work parameter:
//   rwork                      Real array of length n; holds intermediate
//                              solution.
func lusolv(n int, lu []float64, lurow, lcolst, ucolst, rperm, cperm []int, x []float64) error {
	rwork := make([]float64, n)
	err := lsolve(n, lu, lurow, lcolst, ucolst, rperm, cperm, x, rwork)
	if err != nil {
		return err
	}
	err = usolve(n, lu, lurow, lcolst, ucolst, rperm, cperm, rwork, x)
	if err != nil {
		return err
	}
	return nil
}

// cntrow fills its last argument with the nonzero row counts of the
// matrix specified in the first two arguments.
func cntrow(arow []int, lasta int, rowcnt []int) {
	// maxk marks the highest numbered row that has been seen.
	maxk := 0
	for i := 1; i <= lasta; i++ {
		k := arow[i]
		if k > maxk {
			for j := maxk + 1; j <= k; j++ {
				rowcnt[j] = 0
			}
			maxk = k
		}
		rowcnt[k] = rowcnt[k] + 1
	}
}

// rcopy copies a real*8 array A to another array B.
//
// In the following routine for copying whole arrays, the direction
// of iteration (which makes a difference if the arrays overlap) is
// controlled by MODE, which is set false for backward displacement and
// true for forward displacement.
func rcopy(a, b []float64, la int, mode bool) {
	if mode {
		goto l200
	}
	for i := 1; i <= la; i++ {
		b[i] = a[i]
	}
	return

l200:
	for i := la; i >= 1; i-- {
		b[i] = a[i]
	}
	return
}

// icopy copies an integer array A to another array B.
func icopy(a, b []int, la int, mode bool) {
	// In the following routine for copying whole arrays, the direction
	// of iteration (which makes a difference if the arrays overlap) is
	// controlled by mode, which is set false for backward displacement
	// and true for forward displacement.

	if !mode {
		for i := 1; i <= la; i++ {
			b[i] = a[i]
		}
		return
	}
	for i := la; i >= 1; i-- {
		b[i] = a[i]
	}
}

// rfill fills a real*8 array with a given value.
func rfill(a []float64, la int, rval float64) {
	for i := 1; i <= la; i++ {
		a[i] = rval
	}
	return
}

// ifill fills an integer array with a given value.
func ifill(a []int, la, ival int) {
	for i := 1; i <= la; i++ {
		a[i] = ival
	}
}

var rnd int

func dordstat(n, k int, A []float64, kth *float64, info *int) {
	var i, j int
	var x float64

	if k < 0 || k > n {
		*info = -1
		return
	}

	p := 1
	r := n

l100:

	if p == r {
		goto l900
	}

	if r-p >= 8 {
		rnd = 1366*rnd + 150889%714025
		q := p + (rnd % (r - p + 1))

		tmp := A[p]
		A[p] = A[q]
		A[q] = tmp
	}

	x = A[p]
	i = p - 1
	j = r + 1

l200:
	_ = 0

l210:
	j = j - 1
	if A[j] > x {
		goto l210
	}

l220:
	i = i + 1
	if A[i] < x {
		goto l220
	}

	if i < j {
		tmp := A[i]
		A[i] = A[j]
		A[j] = tmp
		goto l200
	}

	if j < k {
		p = j + 1
	} else {
		r = j
	}

	goto l100

l900:
	*kth = A[p]
	*info = 0
	return
}

// requiv tests if two []float64 arrays start at the same address.
func requiv(a, b []float64) bool {
	requiv := false
	temp := a[1]
	a[1] = 0.0
	if b[1] != 0.0 {
		goto l100
	}
	a[1] = 1.0
	if b[1] != 1.0 {
		goto l100
	}
	requiv = true
l100:
	a[1] = temp
	return requiv
}