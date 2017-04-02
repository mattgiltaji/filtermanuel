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
			t.Errorf("ShouldCopy(%v) == %v, want %v", c.in, got, c.want)
		}
	}
	monsterCases := []struct {
		in   string
		want bool
	}{
		//exactly match monsters in list
		{"Monster", true},
		{"Monster 2", true},
		{"Monster'1", true},
		{"comma, the monster", true},
		{"monster-dash", true},
		{"Monster.37", true},

		//simlarly named monsters not in list
		{"monster", false},
		{"yolo", false},
		{"comma, ", false},
		{"-dash", false},
		{"37", false},
		{"Monst", false},

	}
	allowedMonsters := []string{"Monster", "Monster'1", "Monster 2", "Monster.37", "monster-dash", "comma, the monster"}
	for _, c := range monsterCases {
		got := ShouldCopy(c.in, allowedMonsters)
		if got != c.want {
			t.Errorf("ShouldCopy(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}
