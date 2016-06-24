package p5r

// MatchString return true if the string matches the regex
// Returns false if an error/timeout occurs
func (re *Regexp) MatchString(s string) bool {
	m, err := re.run(true, -1, getRunes(s))
	if err != nil {
		return false
	}
	return m != nil
}
