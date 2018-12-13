package readline

import (
	"testing"
)

var runeBufferReplaceRunesTests = []struct {
	Word               string
	Offset             int
	FormatAsIdentifier bool
	AppendSpace        bool
	Buf                string
	Idx                int
	Expect             string
	ExpectIdx          int
}{
	{
		Word:               "replace ",
		Offset:             3,
		FormatAsIdentifier: false,
		AppendSpace:        true,
		Buf:                "abcdef ghi",
		Idx:                10,
		Expect:             "abcdef replace ",
		ExpectIdx:          15,
	},
	{
		Word:               "replace ",
		Offset:             3,
		FormatAsIdentifier: false,
		AppendSpace:        false,
		Buf:                "abcdef ghi",
		Idx:                10,
		Expect:             "abcdef replace",
		ExpectIdx:          14,
	},
	{
		Word:               "replace ",
		Offset:             3,
		FormatAsIdentifier: false,
		AppendSpace:        true,
		Buf:                "abcdef,ghi",
		Idx:                10,
		Expect:             "abcdef, replace ",
		ExpectIdx:          16,
	},
	{
		Word:               "replace ",
		Offset:             0,
		FormatAsIdentifier: false,
		AppendSpace:        true,
		Buf:                "",
		Idx:                0,
		Expect:             "replace ",
		ExpectIdx:          8,
	},
	{
		Word:               "replace ",
		Offset:             3,
		FormatAsIdentifier: false,
		AppendSpace:        true,
		Buf:                "abcdef ghi jkl",
		Idx:                10,
		Expect:             "abcdef replace jkl",
		ExpectIdx:          15,
	},
	{
		Word:               "/path/to/dir/ ",
		Offset:             3,
		FormatAsIdentifier: true,
		AppendSpace:        true,
		Buf:                "abcdef ghi",
		Idx:                10,
		Expect:             "abcdef `/path/to/dir/` ",
		ExpectIdx:          21,
	},
	{
		Word:               "/path/to/dir/ ",
		Offset:             3,
		FormatAsIdentifier: true,
		AppendSpace:        false,
		Buf:                "abcdef ghi",
		Idx:                10,
		Expect:             "abcdef `/path/to/dir/`",
		ExpectIdx:          21,
	},
	{
		Word:               "/path/to/dir/ ",
		Offset:             6,
		FormatAsIdentifier: true,
		AppendSpace:        false,
		Buf:                "abcdef ghi `/path`",
		Idx:                17,
		Expect:             "abcdef ghi `/path/to/dir/`",
		ExpectIdx:          25,
	},
	{
		Word:               "/path/to/dir/ ",
		Offset:             3,
		FormatAsIdentifier: true,
		AppendSpace:        true,
		Buf:                "abcdef ghi",
		Idx:                10,
		Expect:             "abcdef `/path/to/dir/` ",
		ExpectIdx:          21,
	},
	{
		Word:               "column ",
		Offset:             0,
		FormatAsIdentifier: true,
		AppendSpace:        false,
		Buf:                "abcdef `table.csv`.",
		Idx:                19,
		Expect:             "abcdef `table.csv`.column",
		ExpectIdx:          25,
	},
	{
		Word:               "bar() ",
		Offset:             4,
		FormatAsIdentifier: false,
		AppendSpace:        false,
		Buf:                "abcdef foo()",
		Idx:                11,
		Expect:             "abcdef bar()",
		ExpectIdx:          11,
	},
	{
		Word:               "bar() over () ",
		Offset:             4,
		FormatAsIdentifier: false,
		AppendSpace:        false,
		Buf:                "abcdef foo() over ()",
		Idx:                11,
		Expect:             "abcdef bar() over ()",
		ExpectIdx:          11,
	},
}

func TestRuneBuffer_ReplaceRunes(t *testing.T) {
	buf := new(RuneBuffer)
	for _, v := range runeBufferReplaceRunesTests {
		buf.Erase()
		buf.WriteString(v.Buf)
		buf.idx = v.Idx
		buf.ReplaceRunes([]rune(v.Word), v.Offset, v.FormatAsIdentifier, v.AppendSpace)
		result := string(buf.Runes())
		if result != v.Expect || buf.idx != v.ExpectIdx {
			t.Errorf("result = %q, want %q for %q", result, v.Expect, v.Word)
			t.Errorf("index = %d, want %d for %q", buf.idx, v.ExpectIdx, v.Word)
		}
	}
}

var runeBufferFormatAsIdentifierTests = []struct {
	Input        string
	Idx          int
	Expect       string
	ExpectOffset int
}{
	{
		Input:        "ident ",
		Idx:          13,
		Expect:       "ident ",
		ExpectOffset: 0,
	},
	{
		Input:        "ident.ext ",
		Idx:          7,
		Expect:       "`ident.ext` ",
		ExpectOffset: 0,
	},
	{
		Input:        "id`ent ",
		Idx:          7,
		Expect:       "`id\\`ent` ",
		ExpectOffset: 0,
	},
	{
		Input:        "/path/to/file ",
		Idx:          13,
		Expect:       "`/path/to/file` ",
		ExpectOffset: 0,
	},
	{
		Input:        "/path/to/file ",
		Idx:          7,
		Expect:       "`/path/to/file` ",
		ExpectOffset: 0,
	},
	{
		Input:        "/path/to/dir/ ",
		Idx:          13,
		Expect:       "`/path/to/dir/` ",
		ExpectOffset: 2,
	},
}

func TestRuneBuffer_FormatAsIdentifier(t *testing.T) {
	buf := new(RuneBuffer)
	buf.WriteString("`abcde` fghij")
	for _, v := range runeBufferFormatAsIdentifierTests {
		buf.idx = v.Idx
		result, offset := buf.FormatAsIdentifier([]rune(v.Input))
		if string(result) != v.Expect {
			t.Errorf("result = %q, want %q for %q", string(result), v.Expect, v.Input)
		}
		if offset != v.ExpectOffset {
			t.Errorf("offset = %d, want %d for %q", offset, v.ExpectOffset, v.Input)
		}
	}
}
