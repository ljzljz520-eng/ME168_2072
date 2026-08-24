package workflow

import (
	"gymrecommend/internal/model"
	"testing"
)

func TestReport(t *testing.T) {
	r := BuildReport("m", []model.Recommendation{{ClassID: "c", Score: .9, Reason: "fit"}})
	if r.Top.ClassID != "c" || r.Total != 1 {
		t.Fatal(r)
	}
}
