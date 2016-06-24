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

// ReplaceAllString returns a modified string if the replacement worked
// Returns the original string if an error occured
func (re *Regexp) ReplaceAllString(input, replacement string) string {
	output, err := re.Replace(input, replacement, -1, -1)
	if err != nil {
		// Return the original string if something went wrong
		return input
	}
	// Return the string with replacements
	return output
}

func MustCompile(input string) *Regexp {
	return MustCompile2(input, 0)
}
