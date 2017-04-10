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

func filterManuel(manuelPath, faxbotPath, outputPath string) (err error) {
	manuelLines := make(chan string)
	faxbotMonsters := make(chan map[string]struct{})
	firstPassOutput := make(chan string)
	secondPassOutput := make(chan string)

	//todo: read and prep manuel and faxbox in parallel

	go readFaxbotDataFromFile(faxbotPath, faxbotMonsters)
	go readManuelDataFromFile(manuelPath, manuelLines)

	go filterThroughFaxbot(manuelLines, faxbotMonsters, firstPassOutput)
	go removeBlankAreas(firstPassOutput, secondPassOutput)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outputFile.Close()
	output := bufio.NewWriter(outputFile)
	defer output.Flush()

	for line := range secondPassOutput {
		output.WriteString(line)
	}
	return nil

}

func readManuelDataFromFile(manuelPath string, dataChannel chan string) {
	manuelFile, err := os.Open(manuelPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer manuelFile.Close()
	manuel := bufio.NewScanner(manuelFile)
	for manuel.Scan() {
		dataChannel <- manuel.Text()
	}
	close(dataChannel)
}

func readFaxbotDataFromFile(faxbotPath string, dataChannel chan map[string]struct{}) {
	faxbotFile, err := os.Open(faxbotPath)
	if err != nil {
		log.Fatal(err)
	}
	defer faxbotFile.Close()
	faxbot := bufio.NewScanner(faxbotFile)
	faxbotData := make(map[string]struct{})

	for faxbot.Scan() {
		faxbotData[faxbot.Text()] = struct{}{}
	}

	dataChannel <- faxbotData
	close(dataChannel)
}

func filterThroughFaxbot(manuelChannel chan string, faxbotChannel chan map[string]struct{}, outputChannel chan string) {

	faxbotMonsters := <-faxbotChannel
	newLine := ""
	for line := range manuelChannel {
		if shouldCopy(line, faxbotMonsters) {
			outputChannel <- string(newLine + line)
			newLine = "\r\n"
		}
	}
	close(outputChannel)
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

func removeBlankAreas(dataChannel, outputChannel chan string) {
	newLine := ""
	nextLine := ""
	line, ok := <-dataChannel
	for ok {
		if isSectionHeader(line) {
			nextLine, ok = <-dataChannel
			if isSectionSeparator(nextLine) {
				continue
			} else {
				outputChannel <- string(newLine + line)
				newLine = "\r\n"
				outputChannel <- string(newLine + nextLine)
			}
		} else {
			outputChannel <- string(newLine + line)
			newLine = "\r\n"
		}
	}
	close(outputChannel)
}
