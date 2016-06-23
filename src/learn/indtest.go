// Package learn contains the structural learning algorithm as well as a k-means clustering
// and a independence test.
package learn

import "math"

// Lower incomplete Gamma function.
func igamma(a, x float64) float64 {
	var sum float64 = 0.0
	var t float64 = 1.0 / a
	var n float64 = 1.0

	for t != 0 {
		sum += t
		t *= x / (a + n)
		n++
	}

	return math.Pow(x, a) * math.Exp(-x) * sum
}

// Incomplete gamma convergence limit.
//const convgamma = 200

// Incomplete Gamma function.
//func Igamma(k, x float64) float64 {
//if x < 0.0 {
//return 0.0
//}

//s := (1.0 / x) * math.Pow(x, k) * math.Exp(-x)
//sum, nom, den := 1.0, 1.0, 1.0

//for i := 0; i < convgamma; i++ {
//nom *= x
//k++
//den *= k
//sum += nom / den
//}

//return sum * s
//}

//func igammac(a, x float64) float64 {
//if x <= 0 || a <= 0 {
//return 1.0
//} else if x < 1 || x < a {
//return 1.0 - Igamma(a, x)
//}

//lgamma, _ := math.Lgamma(a)
//ax := a*math.Log(x) - x - lgamma
//if ax < -709.78271289338399 {
//return 0.0
//}

//ax = math.Exp(ax)
//var y float64 = 1 - a
//var z float64 = x + y - 1
//c := 0.0
//p2 := 1.0
//q2 := x
//p1 := x + 1
//q1 := z * x
//ans := p1 / q1

//const eps = 0.000000000000001
//const bignum = 4503599627370496.0
//const invbignum = 2.22044604925031308085 * 0.0000000000000001

//var t float64 = -1.0
//var r float64

//for t > eps {
//c++
//y++
//z += 2
//yc := y * c
//pk := p1*z - p2*yc
//qk := q1*z - q2*yc

//if qk <= 0.0 {
//r = pk / qk
//t = math.Abs((ans - r) / r)
//ans = r
//} else {
//t = 1.0
//}

//p2 = p1
//p1 = pk
//q2 = q1
//q1 = qk

//if math.Abs(pk) > bignum {
//p2 = p2 * invbignum
//p1 = p1 * invbignum
//q2 = q2 * invbignum
//q1 = q1 * invbignum
//}
//}

//return ans * ax
//}

//func Igamma(a, x float64) float64 {
//if x <= 0 || a <= 0 {
//return 0.0
//} else if x > 1.0 && x > a {
//return 1.0 - igammac(a, x)
//}

//lgamma, _ := math.Lgamma(a)
//ax := a*math.Log(x) - x - lgamma

//if ax < -709.78271289338399 {
//return 0.0
//}

//ax = math.Exp(ax)
//var r float64 = a
//var c float64 = 1.0
//var ans float64 = 1.0

//const eps = 0.000000000000001

//for c/ans > eps {
//r++
//c = c * x / r
//ans += c
//}

//return ans * ax / a
//}

// Function chisquare returns the p-value of Pr(X^2 > cv).
// Compare this value to the significance level assumed. If chisquare < sigval, then we cannot
// accept the null hypothesis and thus the two variables are dependent.
//
// Thanks to Jacob F. W. for a tutorial on chi-square distributions.
// Source: http://www.codeproject.com/Articles/432194/How-to-Calculate-the-Chi-Squared-P-Value
func Chisquare(df int, cv float64) float64 {
	if cv < 0 || df < 1 {
		return 0.0
	}

	k := float64(df) * 0.5
	x := cv * 0.5

	if df == 2 {
		return math.Exp(-1.0 * x)
	}

	pval := igamma(k, x)

	if math.IsNaN(pval) || math.IsInf(pval, 0) || pval <= 1e-8 {
		return 1e-14
	}

	pval /= math.Gamma(k)

	return 1.0 - pval
}

// ChiSquareTest returns whether variable x and y are statistically independent.
// We use the Chi-Square test to find correlations between the two variables.
// Argument data is a table with the counting of each variable category, where the first axis is
// the counting of each category of variable x and the second axis of variable y. The last element
// of each row and column is the total counting. E.g.:
//
// +------------------------+
// |      X_1 X_2 X_3 total |
// | Y_1  100 200 100  400  |
// | Y_2   50 300  25  375  |
// |total 150 500 125  775  |
// +------------------------+
//
// Argument p is the number of categories (or levels) in x
// Argument q is the number of categories (or levels) in y
//
// Returns true if independent and false otherwise.
func ChiSquareTest(p, q int, data [][]float64) bool {

	// df is the degree of freedom.
	df := (p - 1) * (q - 1)

	// Expected frequencies.
	E := make([][]float64, p)
	for i := 0; i < p; i++ {
		E[i] = make([]float64, q)
	}

	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			E[i][j] = float64(data[i][p]*data[j][q]) / float64(data[p][q])
		}
	}

	// Test statistic.
	var chi float64 = 0
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			diff := float64(data[i][j] - E[i][j])
			chi += (diff * diff) / E[i][j]
		}
	}

	// Significance value.
	const sigval = 0.05

	// Compare cmd with sigval. If cmp < sigval, then dependent. Otherwise independent.
	cmp := Chisquare(df, chi)

	return cmp >= sigval
}
