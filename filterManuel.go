package main

import "regexp"

var (
	sectionSeparator = regexp.MustCompile(`^=====*$`)
	sectionHeader = regexp.MustCompile(`^\[[A-Z]+.*\]$`)
)

func ShouldCopy(s string, allowed map[string]struct{}) bool {
	if sectionSeparator.MatchString(s) || sectionHeader.MatchString(s) {
		return true
	}
	_, ok := allowed[s]
	return ok
}
