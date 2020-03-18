package sha

import "testing"

func TestSha256(t *testing.T) {
	file := "origin.txt"
	sha := GetSha256(file)
	if "46cbee948bb7bcd507817a50beb2381ec8983a2e7d96aef1636e6a6a3328634d" != sha {
		t.Error("sha output not equal")
	}
}
