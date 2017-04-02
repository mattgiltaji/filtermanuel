package main

import "regexp"

var (
	sectionSeparator = regexp.MustCompile(`^=====*$`)
	sectionHeader = regexp.MustCompile(`^\[[A-Z]+.*\]$`)
)

func ShouldCopy(s string, allowed []string) bool {
	_ = allowed
	return sectionSeparator.MatchString(s) || sectionHeader.MatchString(s)
}
