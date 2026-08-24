package model

import "testing"

func TestMemberValidation(t *testing.T) {
	if ValidateMember(Member{ID: "m", Name: "A", CoachRating: 4}) != nil {
		t.Fatal("valid member rejected")
	}
	if ValidateMember(Member{ID: "", Name: "A"}) == nil {
		t.Fatal("missing id accepted")
	}
}
