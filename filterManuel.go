package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

var (
	sectionSeparator = regexp.MustCompile(`^=====*$`)
	sectionHeader    = regexp.MustCompile(`^\[[A-Z]+.*\]$`)
	bracketsToIgnore = regexp.MustCompile(` \{[0-3]\}$`)
)

func filterManuel(manuelPath, faxbotPath, outputPath string) {
	//todo: read and prep manuel and faxbox in parallel
	manuelFile, err := os.Open(manuelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer manuelFile.Close()
	manuel := bufio.NewScanner(manuelFile)

	faxbotFile, err := os.Open(faxbotPath)
	if err != nil {
		log.Fatal(err)
	}
	defer faxbotFile.Close()
	faxbot := bufio.NewScanner(faxbotFile)
	faxbotMonsters := make(map[string]struct{})

	for faxbot.Scan() {
		faxbotMonsters[faxbot.Text()] = struct{}{}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	output := bufio.NewWriter(outputFile)
	defer output.Flush()
	newLine := ""
	for manuel.Scan() {
		line := manuel.Text()
		if shouldCopy(line, faxbotMonsters) {
			_, err := output.WriteString(newLine + line)
			if err != nil {
				log.Fatal(err)
			}
			newLine = "\r\n"
		}
	}
}

func shouldCopy(s string, allowed map[string]struct{}) bool {
	if isSectionSeparator(s) {
		return true
	}
	if isSectionHeader(s) {
		return true
	}
	if isInAllowedSet(s, allowed) {
		return true
	}
	return false
}

func isSectionSeparator(s string) bool {
	return sectionSeparator.MatchString(s)
}

func isSectionHeader(s string) bool {
	return sectionHeader.MatchString(s)
}

func isInAllowedSet(s string, allowed map[string]struct{}) bool {
	stringWithBracketsRemoved := bracketsToIgnore.ReplaceAllString(s, "")
	_, ok := allowed[stringWithBracketsRemoved]
	return ok
}
