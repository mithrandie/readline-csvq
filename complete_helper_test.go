package readline

import (
	"testing"
)

var lastElementTests = []struct {
	Input  string
	Expect string
}{
	{
		Input:  "abcdef ghij",
		Expect: "ghij",
	},
	{
		Input:  "abcdef ghij ",
		Expect: "",
	},
	{
		Input:  "abcdef `gh\\`ij",
		Expect: "`gh`ij",
	},
	{
		Input:  "abcdef `gh`ij",
		Expect: "ij",
	},
	{
		Input:  "abcdef 'gh\\'ij",
		Expect: "'gh'ij",
	},
	{
		Input:  "abcdef \"gh\\\"ij",
		Expect: "\"gh\"ij",
	},
	{
		Input:  "abcdef @%`ENV`",
		Expect: "@%`ENV`",
	},
	{
		Input:  "abcdef ghij,",
		Expect: ",",
	},
	{
		Input:  "abcdef,ghij",
		Expect: "ghij",
	},
	{
		Input:  "abcdef  ghij",
		Expect: "ghij",
	},
	{
		Input:  "abcdef ;ghij",
		Expect: "ghij",
	},
	{
		Input:  "a,",
		Expect: ",",
	},
	{
		Input:  "abc/def",
		Expect: "abc/def",
	},
	{
		Input:  "`/abc/def",
		Expect: "`/abc/def",
	},
	{
		Input:  "abcdefghij;",
		Expect: "",
	},
}

func TestLastElement(t *testing.T) {
	for _, v := range lastElementTests {
		result := LastElement(v.Input)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q", result, v.Expect, v.Input)
		}
	}
}

var literalIsEnclosedTests = []struct {
	Input  string
	Mark   rune
	Expect bool
}{
	{
		Input:  "abcdef",
		Mark:   '"',
		Expect: true,
	},
	{
		Input:  "abc'defghi'jkl",
		Mark:   '\'',
		Expect: true,
	},
	{
		Input:  "abc'def\\'ghi'jkl",
		Mark:   '\'',
		Expect: true,
	},
	{
		Input:  "abc'defghijkl",
		Mark:   '\'',
		Expect: false,
	},
}

func TestLiteralIsEnclosed(t *testing.T) {
	for _, v := range literalIsEnclosedTests {
		result := LiteralIsEnclosed(v.Mark, []rune(v.Input))
		if result != v.Expect {
			t.Errorf("result = %t, want %t for %q", result, v.Expect, v.Input)
		}
	}
}

var bracketlIsEnclosedTests = []struct {
	Input  string
	Mark   rune
	Expect bool
}{
	{
		Input:  "abc(defghi)jkl",
		Mark:   '(',
		Expect: true,
	},
	{
		Input:  "abc(defghi\\)jkl",
		Mark:   '(',
		Expect: false,
	},
	{
		Input:  "abc(()defghijkl",
		Mark:   '(',
		Expect: false,
	},
}

func TestBracketIsEnclosed(t *testing.T) {
	for _, v := range bracketlIsEnclosedTests {
		result := BracketIsEnclosed(v.Mark, []rune(v.Input))
		if result != v.Expect {
			t.Errorf("result = %t, want %t for %q", result, v.Expect, v.Input)
		}
	}
}
