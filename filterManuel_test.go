package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFilterManuel(t *testing.T) {
	testDataDir := getTestDataDir(t)
	cases := []string{
		filepath.Join(testDataDir, "copyEverything"),
		filepath.Join(testDataDir, "copyNothing"),
		filepath.Join(testDataDir, "copySomething"),
		filepath.Join(testDataDir, "fullTest"),
	}
	for _, c := range cases {
		gotFile, err := ioutil.TempFile("", "filtered_manuel")
		if err != nil {
			t.Fatalf("Could not create temp output file. Error: %v", err)
		}
		defer os.Remove(gotFile.Name())

		manuelFile := filepath.Join(c, "manuel.txt")
		faxbotFile := filepath.Join(c, "faxbot.txt")
		expectedFile := filepath.Join(c, "expected.txt")

		filterManuel(manuelFile, faxbotFile, gotFile.Name())

		expectedContents, err := ioutil.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("Could not read %v. Error: %v", expectedFile, err)
		}
		gotContents, err := ioutil.ReadFile(gotFile.Name())
		if err != nil {
			t.Fatalf("Could not read %v. Error: %v", gotFile.Name(), err)
		}
		defer gotFile.Close()
		if string(expectedContents) != string(gotContents) {
			t.Errorf("filterManuel(%v, %v, %v), == '%v', want '%v'", manuelFile, faxbotFile,
				gotFile.Name(), string(gotContents), string(expectedContents))
		}

	}
}

func TestShouldCopy(t *testing.T) {
	allowedMonsters := map[string]struct{}{
		"Monster": {}, "Monster'1": {}, "Monster 2": {}, "Monster.37": {},
		"monster-dash": {}, "comma, the monster": {}}
	cases := []struct {
		in      string
		allowed map[string]struct{}
		want    bool
	}{
		//good section separators
		{"=================================", nil, true},     //real line separator
		{"=====================================", nil, true}, //longer
		{"==========", nil, true},                            //arbitrarily short but valid
		{"====", nil, true},                                  //shortest possible valid
		//bad section separators
		{"=", nil, false},                             //too short
		{"==", nil, false},                            //also too short
		{"===", nil, false},                           //still too short
		{"---------------------------", nil, false},   //wrong symbol
		{"-=-=-=-=-=-=-=-=-=-=", nil, false},          //no mixing
		{"----=======-----", nil, false},              //nope
		{"============================-", nil, false}, //no dash

		//good area names
		{"[Ye Olde Medievale Villagee]", nil, true},                  // plain name
		{"[An Incredibly Strange Place (Mediocre Trip)]", nil, true}, //parentheses
		{"[Anger Man's Level]", nil, true},                           //apostrophe
		{"[The Gourd!]", nil, true},                                  //bang!
		{"[LavaCo™ Lamp Factory]", nil, true},                        //trademark
		{"[A Deserted Stretch of I-911]", nil, true},                 //dash
		{"[Engineering]", nil, true},                                 // no space in area name
		//bad area names
		{"[]", nil, false}, {"[ ]", nil, false}, //blanks alone insufficient
		{"[!]", nil, false}, {"[*****]", nil, false}, //alphabetical rune needed

		//exactly match monsters in list
		{"Monster", allowedMonsters, true},
		{"Monster 2", allowedMonsters, true},
		{"Monster'1", allowedMonsters, true},
		{"comma, the monster", allowedMonsters, true},
		{"monster-dash", allowedMonsters, true},
		{"Monster.37", allowedMonsters, true},

		//simlarly named monsters not in list
		{"monster", allowedMonsters, false},
		{"yolo", allowedMonsters, false},
		{"comma, ", allowedMonsters, false},
		{"-dash", allowedMonsters, false},
		{"37", allowedMonsters, false},
		{"Monst", allowedMonsters, false},

		//ignore the braces after monster name that missingManuel adds
		{"Monster {1}", allowedMonsters, true},
		{"monster-dash {2}", allowedMonsters, true},
		{"Monster.37 {3}", allowedMonsters, true},
		{"Monst {3}", allowedMonsters, false},
	}
	for _, c := range cases {
		got := shouldCopy(c.in, c.allowed)
		if got != c.want {
			t.Errorf("shouldCopy(%v, %v) == %v, want %v", c.in, c.allowed, got, c.want)
		}
	}
}

func TestRemoveBlankAreas(t *testing.T) {
	testDataDir := getTestDataDir(t)
	cases := []string{
		filepath.Join(testDataDir, "removeEverything"),
		filepath.Join(testDataDir, "removeNothing"),
		filepath.Join(testDataDir, "removeSomething"),
	}
	for _, c := range cases {
		toFilterFile := filepath.Join(c, "contents.txt")
		wantFile := filepath.Join(c, "wanted.txt")

		toFilterContents, err := ioutil.ReadFile(toFilterFile)
		if err != nil {
			t.Fatalf("Could not read %v. Error: %v", toFilterFile, err)
		}
		wantContents, err := ioutil.ReadFile(wantFile)
		if err != nil {
			t.Fatalf("Could not read %v. Error: %v", wantFile, err)
		}
		toFilter := bufio.NewScanner(bytes.NewReader(toFilterContents))
		want := bufio.NewScanner(bytes.NewReader(wantContents))

		gotRaw := removeBlankAreas(toFilter)
		got := bufio.NewScanner(strings.NewReader(strings.Join(gotRaw, "")))
		i := 1
		for want.Scan() && got.Scan() {
			wantLine := want.Text()
			gotLine := got.Text()
			if wantLine != gotLine {
				t.Errorf("line %v of removeBlankAreas(%v) == '%v', want '%v'", i, toFilterFile, gotLine, wantLine)
				break
			}
			i++
		}

	}
}

func getTestDataDir(t *testing.T) (testDataDir string) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Unable to determine runtime file location.")
	}
	thisDir := filepath.Dir(thisFile)
	testDataDir, err := filepath.Abs(filepath.Join(thisDir, "testdata"))
	if err != nil {
		t.Fatalf("Unable to find testdata directory. Err: %v", err)
	}
	return
}
