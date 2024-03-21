package base

import "testing"

func TestFormatCode(t *testing.T) {

	shells := [][]string{
		{"go", "mod", "tidy"},
		{"gofmt", "-w", "."},
	}

	if err := FormatCode("./", shells...); err != nil {
		t.Error(err)
	}

}
