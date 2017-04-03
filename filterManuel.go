package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	sectionSeparator = regexp.MustCompile(`^=====*$`)
	sectionHeader    = regexp.MustCompile(`^\[[A-Z]+.*\]$`)
	bracketsToIgnore = regexp.MustCompile(` \{[0-3]\}$`)
	manuelFilePath   string
	faxbotFilePath   string
	outputFilePath   string
)

func init() {
	const (
		manuelUsage = "path to the missingManuel.ash data file"
		faxbotUsage = "path to the faxbot data file"
		outputUsage = "path to where the output data should be written"
	)
	var (
		defaultDir    = filepath.Clean(`C:\Users\admin\Desktop\kolmafia\samples`)
		manuelDefault = filepath.Join(defaultDir, "monster manuel.txt")
		faxbotDefault = filepath.Join(defaultDir, "faxbot.txt")
		outputDefault = filepath.Join(defaultDir, "filtered_faxbot.txt")
	)

	flag.StringVar(&manuelFilePath, "manuel", manuelDefault, manuelUsage)
	flag.StringVar(&manuelFilePath, "m", manuelDefault, manuelUsage+" (shorthand)")
	flag.StringVar(&faxbotFilePath, "faxbot", faxbotDefault, faxbotUsage)
	flag.StringVar(&faxbotFilePath, "f", faxbotDefault, faxbotUsage+" (shorthand)")
	flag.StringVar(&outputFilePath, "output", outputDefault, outputUsage)
	flag.StringVar(&outputFilePath, "o", outputDefault, outputUsage+" (shorthand)")
}

func main() {
	flag.Parse()
	filterManuel(manuelFilePath, faxbotFilePath, outputFilePath)
}

func filterManuel(manuelPath, faxbotPath, outputPath string) {
	//todo: read and prep manuel and faxbox in parallel
	manuelFile, err := os.Open(manuelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer manuelFile.Close()
	manuel := bufio.NewScanner(manuelFile)

	faxbotMonsters := getFaxbotData(faxbotPath)

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
	//read from output, filter through removeBlankAreas() and write it out again
}

func getFaxbotData(faxbotPath string) (faxbotData map[string]struct{}) {
	faxbotFile, err := os.Open(faxbotPath)
	if err != nil {
		log.Fatal(err)
	}
	defer faxbotFile.Close()
	faxbot := bufio.NewScanner(faxbotFile)
	faxbotData = make(map[string]struct{})

	for faxbot.Scan() {
		faxbotData[faxbot.Text()] = struct{}{}
	}
	return
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

func removeBlankAreas(contents *bufio.Scanner) (filtered []string) {
	newLine := ""
	for contents.Scan() {
		line := contents.Text()
		if isSectionHeader(line) {
			contents.Scan()
			nextLine := contents.Text()
			if isSectionSeparator(nextLine) {
				continue
			} else {
				filtered = append(filtered, newLine+line)
				newLine = "\r\n"
				filtered = append(filtered, newLine+nextLine)
			}
		} else {
			filtered = append(filtered, newLine+line)
			newLine = "\r\n"
		}
	}
	return
}
