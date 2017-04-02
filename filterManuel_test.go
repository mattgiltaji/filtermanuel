package main

import "testing"

func TestShouldCopy(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		//good section separators
		{"=================================", true},     //real line separator
		{"=====================================", true}, //longer
		{"==========", true},                            //arbitrarily short but valid
		{"====", true},                                  //shortest possible valid
		//bad section separators
		{"=", false},                             //too short
		{"==", false},                            //also too short
		{"===", false},                           //still too short
		{"---------------------------", false},   //wrong symbol
		{"-=-=-=-=-=-=-=-=-=-=", false},          //no mixing
		{"----=======-----", false},              //nope
		{"============================-", false}, //no dash

		//good area names
		{"[Ye Olde Medievale Villagee]", true},                  // plain name
		{"[An Incredibly Strange Place (Mediocre Trip)]", true}, //parentheses
		{"[Anger Man's Level]", true},                           //apostrophe
		{"[The Gourd!]", true},                                  //bang!
		{"[LavaCo™ Lamp Factory]", true},                        //trademark
		{"[A Deserted Stretch of I-911]", true},                 //dash
		{"[Engineering]", true},                                 // no space in area name
		//bad area names
		{"[]", false}, {"[ ]", false},      //blanks alone insufficient
		{"[!]", false}, {"[*****]", false}, //alphabetical rune needed

	}
	for _, c := range cases {
		got := ShouldCopy(c.in, nil)
		if got != c.want {
			t.Error("ShouldCopy(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
