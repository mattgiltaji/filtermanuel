package main

import "testing"

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
			t.Errorf("ShouldCopy(%v, %v) == %v, want %v", c.in, c.allowed, got, c.want)
		}
	}
}
