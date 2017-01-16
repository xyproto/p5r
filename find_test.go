package p5r

import "testing"

// Please note that regex for Go and Javascript/Perl is slightly different
// I.e pattern `a*` for string `baaab` expects (0,0), (1,4), (5, 5) for Go
// and for Javascript/Perl it additionally expects (4,4) - this version is used in code.
// Moreover p5r does not support non-english characters as it returns incorrect results -
// this tests have been commented out.

// For each pattern/text pair, what is the expected output of each function?
// We can derive the textual results from the indexed results, the non-submatch
// results from the submatched results, the single results from the 'all' results,
// and the byte results from the string results. Therefore the table includes
// only the FindAllStringSubmatchIndex result.
type FindTest struct {
	pat     string
	text    string
	matches [][]int
}

var findTests = []FindTest{
	{``, ``, build(1, 0, 0)},
	{`^abcdefg`, "abcdefg", build(1, 0, 7)},
	{`a+`, "baaab", build(1, 1, 4)},
	{"abcd..", "abcdef", build(1, 0, 6)},
	{`a`, "a", build(1, 0, 1)},
	{`x`, "y", nil},
	{`b`, "abc", build(1, 1, 2)},
	{`.`, "a", build(1, 0, 1)},
	{`.*`, "abcdef", build(1, 0, 6)},
	{`^`, "abcde", build(1, 0, 0)},
	{`$`, "abcde", build(1, 5, 5)},
	{`^abcd$`, "abcd", build(1, 0, 4)},
	{`^bcd'`, "abcdef", nil},
	{`^abcd$`, "abcde", nil},
	{`a+`, "baaab", build(1, 1, 4)},
	{`a*`, "baaab", build(4, 0, 0, 1, 4, 4, 4, 5, 5)},
	{`[a-z]+`, "abcd", build(1, 0, 4)},
	{`[^a-z]+`, "ab1234cd", build(1, 2, 6)},
	{`[a\-\]z]+`, "az]-bcz", build(2, 0, 4, 6, 7)},
	{`[^\n]+`, "abcd\n", build(1, 0, 4)},
	{`()`, "", build(1, 0, 0, 0, 0)},
	{`(a)`, "a", build(1, 0, 1, 0, 1)},
	{`(.)(.)`, "ba", build(1, 0, 2, 0, 1, 1, 2)},
	//{`[日本語]+`, "日本語日本語", build(1, 0, 18)},
	//{`日本語+`, "日本語", build(1, 0, 9)},
	//{`日本語+`, "日本語語語語", build(1, 0, 18)},
	//{`(.)(.)`, "日a", build(1, 0, 4, 0, 3, 3, 4)},
	{`(.*)`, "", build(1, 0, 0, 0, 0)},
	{`(.*)`, "abcd", build(1, 0, 4, 0, 4)},
	{`(..)(..)`, "abcd", build(1, 0, 4, 0, 2, 2, 4)},
	{`(([^xyz]*)(d))`, "abcd", build(1, 0, 4, 0, 4, 0, 3, 3, 4)},
	{`((a|b|c)*(d))`, "abcd", build(1, 0, 4, 0, 4, 2, 3, 3, 4)},
	{`(((a|b|c)*)(d))`, "abcd", build(1, 0, 4, 0, 4, 0, 3, 2, 3, 3, 4)},
	{`\a\f\n\r\t\v`, "\a\f\n\r\t\v", build(1, 0, 6)},
	{`[\a\f\n\r\t\v]+`, "\a\f\n\r\t\v", build(1, 0, 6)},

	{`a*(|(b))c*`, "aacc", build(1, 0, 4, 2, 2, -1, -1)},
	{`(.*).*`, "ab", build(1, 0, 2, 0, 2)},
	{`[.]`, ".", build(1, 0, 1)},
	{`/$`, "/abc/", build(1, 4, 5)},
	{`/$`, "/abc", nil},

	// multiple matches
	{`.`, "abc", build(3, 0, 1, 1, 2, 2, 3)},
	{`(.)`, "abc", build(3, 0, 1, 0, 1, 1, 2, 1, 2, 2, 3, 2, 3)},
	{`.(.)`, "abcd", build(2, 0, 2, 1, 2, 2, 4, 3, 4)},
	{`ab*`, "abbaab", build(3, 0, 3, 3, 4, 4, 6)},
	{`a(b*)`, "abbaab", build(3, 0, 3, 1, 3, 3, 4, 4, 4, 4, 6, 5, 6)},

	// fixed bugs
	{`ab$`, "cab", build(1, 1, 3)},
	{`axxb$`, "axxcb", nil},
	{`data`, "daXY data", build(1, 5, 9)},
	{`da(.)a$`, "daXY data", build(1, 5, 9, 7, 8)},
	{`zx+`, "zzx", build(1, 1, 3)},
	{`ab$`, "abcab", build(1, 3, 5)},
	{`(aa)*$`, "a", build(1, 1, 1, -1, -1)},
	{`(?:.|(?:.a))`, "", nil},
	{`(?:A(?:A|a))`, "Aa", build(1, 0, 2)},
	{`(?:A|(?:A|a))`, "a", build(1, 0, 1)},
	{`(a){0}`, "", build(1, 0, 0, -1, -1)},
	{`(?-s)(?:(?:^).)`, "\n", nil},
	{`(?s)(?:(?:^).)`, "\n", build(1, 0, 1)},
	{`(?:(?:^).)`, "\n", nil},
	{`\b`, "x", build(2, 0, 0, 1, 1)},
	{`\b`, "xx", build(2, 0, 0, 2, 2)},
	{`\b`, "x y", build(4, 0, 0, 1, 1, 2, 2, 3, 3)},
	{`\b`, "xx yy", build(4, 0, 0, 2, 2, 3, 3, 5, 5)},
	{`\B`, "x", nil},
	{`\B`, "xx", build(1, 1, 1)},
	{`\B`, "x y", nil},
	{`\B`, "xx yy", build(2, 1, 1, 4, 4)},

	// RE2 tests
	{`[^\S\s]`, "abcd", nil},
	{`[^\S[:space:]]`, "abcd", nil},
	{`[^\D\d]`, "abcd", nil},
	{`[^\D[:digit:]]`, "abcd", nil},
	{`(?i)\W`, "x", nil},
	{`(?i)\W`, "k", nil},
	{`(?i)\W`, "s", nil},

	// can backslash-escape any punctuation
	//{`\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~`,
	// `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`, build(1, 0, 31)},
	//{`[\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~]+`,
	// `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`, build(1, 0, 31)},
	{"\\`", "`", build(1, 0, 1)},
	{"[\\`]+", "`", build(1, 0, 1)},

	// long set of matches (longer than startSize)
	{
		".",
		"qwertyuiopasdfghjklzxcvbnm1234567890",
		build(36, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10,
			10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 20,
			20, 21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 29, 29, 30,
			30, 31, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36),
	},
}

func TestFindAllStringIndex(t *testing.T) {
	for _, test := range findTests {
		testFindAllIndex(&test, MustCompile(test.pat).FindAllStringIndex(test.text, -1), t)
	}
}

func TestFindAllSubmatchIndex(t *testing.T) {
	for _, test := range findTests {
		testFindAllSubmatchIndex(&test, MustCompile(test.pat).FindAllSubmatchIndex([]byte(test.text), -1), t)
	}
}

func TestFindAllStringSubmatchIndex(t *testing.T) {
	for _, test := range findTests {
		testFindAllSubmatchIndex(&test, MustCompile(test.pat).FindAllStringSubmatchIndex(test.text, -1), t)
	}
}

func TestFindStringSubmatchIndex(t *testing.T) {
	for _, test := range findTests {
		testFindSubmatchIndex(&test, MustCompile(test.pat).FindStringSubmatchIndex(test.text), t)
	}
}

func TestFindStringIndex(t *testing.T) {
	for _, test := range findTests {
		testFindIndex(&test, MustCompile(test.pat).FindStringIndex(test.text), t)
	}
}

// build is a helper to construct a [][]int by extracting n sequences from x.
// This represents n matches with len(x)/n submatches each.
func build(n int, x ...int) [][]int {
	ret := make([][]int, n)
	runLength := len(x) / n
	j := 0
	for i := range ret {
		ret[i] = make([]int, runLength)
		copy(ret[i], x[j:])
		j += runLength
		if j > len(x) {
			panic("invalid build entry")
		}
	}
	return ret
}

func testFindAllIndex(test *FindTest, result [][]int, t *testing.T) {
	switch {
	case test.matches == nil && result == nil:
	// ok
	case test.matches == nil && result != nil:
		t.Errorf("%s: expected no match; got one: %s", test.pat, test)
	case test.matches != nil && result == nil:
		t.Errorf("%s: expected match; got none: %s", test.pat, test)
	case test.matches != nil && result != nil:
		if len(test.matches) != len(result) {
			t.Errorf("%s: expected %d matches; got %d: %s", test.pat, len(test.matches), len(result), test)
			return
		}
		for k, e := range test.matches {
			if e[0] != result[k][0] || e[1] != result[k][1] {
				t.Errorf("%s: match %d: expected %v got %v: %s", test.pat, k, e, result[k], test)
			}
		}
	}
}

func testFindSubmatchIndex(test *FindTest, result []int, t *testing.T) {
	switch {
	case test.matches == nil && result == nil:
	// ok
	case test.matches == nil && result != nil:
		t.Errorf("expected no match; got one: %s", test)
	case test.matches != nil && result == nil:
		t.Errorf("expected match; got none: %s", test)
	case test.matches != nil && result != nil:
		testSubmatchIndices(test, 0, test.matches[0], result, t)
	}
}

func testFindAllSubmatchIndex(test *FindTest, result [][]int, t *testing.T) {
	switch {
	case test.matches == nil && result == nil:
	// ok
	case test.matches == nil && result != nil:
		t.Errorf("%s: expected no match; got one: %s", test.pat, test)
	case test.matches != nil && result == nil:
		t.Errorf("%s: expected match; got none: %s", test.pat, test)
	case len(test.matches) != len(result):
		t.Errorf("%s: expected %d matches; got %d: %s", test.pat, len(test.matches), len(result), test)
	case test.matches != nil && result != nil:
		for k, match := range test.matches {
			testSubmatchIndices(test, k, match, result[k], t)
		}
	}
}

func testSubmatchIndices(test *FindTest, n int, expect, result []int, t *testing.T) {
	if len(expect) != len(result) {
		t.Errorf("%s: match %d: expected %d matches; got %d: %s", test.pat, n, len(expect)/2, len(result)/2, test)
		return
	}
	for k, e := range expect {
		if e != result[k] {
			t.Errorf("%s: match %d: submatch error: expected %v got %v: %s", test.pat, n, expect, result, test)
		}
	}
}

func testFindIndex(test *FindTest, result []int, t *testing.T) {
	switch {
	case len(test.matches) == 0 && len(result) == 0:
	// ok
	case test.matches == nil && result != nil:
		t.Errorf("expected no match; got one: %s", test)
	case test.matches != nil && result == nil:
		t.Errorf("expected match; got none: %s", test)
	case test.matches != nil && result != nil:
		expect := test.matches[0]
		if expect[0] != result[0] || expect[1] != result[1] {
			t.Errorf("expected %v got %v: %s", expect, result, test)
		}
	}
}
