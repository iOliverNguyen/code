package ss_test

import (
	"reflect"
	"testing"

	. "github.com/ng-vu/go-stuff/sudoku/ss"
)

var input string = `
  000 507 000
  002 406 300
  090 010 020

  270 000 068
  003 000 100
  140 000 093

  060 040 050
  009 205 600
  000 903 000`

func TestSudoku(T *testing.T) {
	expected := Board{
		{0, 0, 0, 5, 0, 7, 0, 0, 0},
		{0, 0, 2, 4, 0, 6, 3, 0, 0},
		{0, 9, 0, 0, 1, 0, 0, 2, 0},

		{2, 7, 0, 0, 0, 0, 0, 6, 8},
		{0, 0, 3, 0, 0, 0, 1, 0, 0},
		{1, 4, 0, 0, 0, 0, 0, 9, 3},

		{0, 6, 0, 0, 4, 0, 0, 5, 0},
		{0, 0, 9, 2, 0, 5, 6, 0, 0},
		{0, 0, 0, 9, 0, 3, 0, 0, 0},
	}

	b, err := Parse(input)
	if err != nil {
		T.Errorf("Unable to parse: %v\n%v", err, b)
		return
	}

	if !reflect.DeepEqual(b, expected) {
		T.Errorf("Unexpected board:\n%v\n\n%v", b, expected)
		return
	}

	err = IsValid(b)
	if err != nil {
		T.Error("Expected valid:", err)
		return
	}

	{
		// Test IsValid
		b1 := b
		b1[0][7] = 7
		err := IsValid(b1)
		if err == nil {
			T.Error("Expected invalid")
		}

		b2 := b
		b2[5][5] = 4
		err = IsValid(b2)
		if err == nil {
			T.Error("Expected invalid")
		}

		b3 := b
		b3[4][8] = 6
		err = IsValid(b3)
		if err == nil {
			T.Error("Expected invalid")
		}
	}

	{
		// Test FindOne
		result, ok := FindOne(b)
		if !ok {
			T.Error("Expected to find at least a result")
			return
		}

		err = IsCompleted(result)
		if err != nil {
			T.Errorf("Expected result:\n%v", result)
			return
		}

		T.Logf("Result:\n%v", result)
	}
}
