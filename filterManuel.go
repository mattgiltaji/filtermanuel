package main

import "regexp"

var (
	sectionSeparator = regexp.MustCompile(`^=====*$`)
	sectionHeader    = regexp.MustCompile(`^\[[A-Z]+.*\]$`)
	bracketsToIgnore = regexp.MustCompile(` \{[0-3]\}$`)
)

func ShouldCopy(s string, allowed map[string]struct{}) bool {
	/*
	Returns true if string should be copied to output file

	3 types of lines should return true:
	    A) a monster that matches a line in allowed (ignoring special brackets)
	    B) a section header - [foo bar]
	    C) a section divider - =====...===
	We don't necessarily evaluate these in order, we try the easy ones first
	*/
	if sectionSeparator.MatchString(s) || sectionHeader.MatchString(s) {
		return true
	}
	_, ok := allowed[bracketsToIgnore.ReplaceAllString(s, "")]
	return ok
}
