package common

import "testing"

func TestUrlAnalysis(t *testing.T) {
	url := "/v1/test/1"
	out := UrlAnalysis(url)
	tmp := []string{"v1", "test", "1"}
	if len(out) != len(tmp) {
		t.Error("Split error")
		t.FailNow()
	}
	for i := 0; i < len(out); i++ {
		if out[i] != tmp[i] {
			t.Errorf("data error %v = %v", out[i], tmp[i])
			t.FailNow()
		}
	}
	t.Log("----------------- urlAnalysis OK --------------------------")
}
