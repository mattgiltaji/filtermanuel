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
	/*Returns true if string should be copied to output file

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
