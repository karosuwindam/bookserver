package writetable

import (
	"fmt"
	"testing"
)

func TestAnalys(t *testing.T) {
	filename := "(aa)[bb]cc.zip"
	a, n := removeParentheses(filename)
	if len(a) == 0 || a[0] != "aa" || n != "[bb]cc.zip" {
		t.Fatalf("%v,%v", fmt.Sprintln(a), n)
	}
	a1, n1 := removeBrackets(filename)
	if len(a1) == 0 || a1[0] != "bb" || n1 != "(aa)cc.zip" {
		t.Fatalf("%v,%v", fmt.Sprintln(a1), n1)
	}
	filename = "(aa)(bb)[cc][dd]ee.zip"
	a, n = removeParentheses(filename)
	if len(a) == 0 || a[0] != "aa" || a[1] != "bb" || n != "[cc][dd]ee.zip" {
		t.Fatalf("%v,%v", fmt.Sprintln(a), n)
	}
	a1, n1 = removeBrackets(filename)
	if len(a1) == 0 || a1[0] != "cc" || a1[1] != "dd" || n1 != "(aa)(bb)ee.zip" {
		t.Fatalf("%v,%v", fmt.Sprintln(a1), n1)
	}
}
