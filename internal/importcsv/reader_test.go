package importcsv

import (
	"strings"
	"testing"
)

func TestReadMembers(t *testing.T) {
	xs, e := ReadMembers(strings.NewReader("member_id,name,goals,slot,coach_rating,package\nm,A,strength,morning,4.5,group"))
	if e != nil || len(xs) != 1 {
		t.Fatal(xs, e)
	}
}
