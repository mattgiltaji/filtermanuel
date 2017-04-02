package main

import "testing"

func TestShouldCopy(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"=================================", true},     //real line separator
		{"=====================================", true}, //longer
		{"==========", true},                            //arbitrarily short but valid
		{"====", true},                                  //shortest possible valid
		{"=", false},                                    //too short
		{"==", false},                                   //also too short
		{"===", false},                                  //still too short
		{"---------------------------", false},          //wrong symbol
		{"-=-=-=-=-=-=-=-=-=-=", false},                 //no mixing
		{"----=======-----", false},                     //nope
		{"============================-", false},        //no dash

	}
	for _, c := range cases {
		got := ShouldCopy(c.in, nil)
		if got != c.want {
			t.Error("ShouldCopy(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
