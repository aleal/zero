package metadata

import "testing"

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestGetVersionConsistent(t *testing.T) {
	v1 := GetVersion()
	v2 := GetVersion()
	if v1 != v2 {
		t.Errorf("GetVersion() not consistent: %q vs %q", v1, v2)
	}
}
