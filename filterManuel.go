package main

import "regexp"

var sectionSeparator = regexp.MustCompile(`^=====*$`)

func ShouldCopy(s string, allowed []string) bool {
	_ = allowed
	return sectionSeparator.MatchString(s)
}
