package readline

import (
	"reflect"
	"testing"
)

func rs(s [][]rune) []string {
	ret := make([]string, len(s))
	for idx, ss := range s {
		ret[idx] = string(ss)
	}
	return ret
}

func sr(s ...string) [][]rune {
	ret := make([][]rune, len(s))
	for idx, ss := range s {
		ret[idx] = []rune(ss)
	}
	return ret
}

func TestRetSegment(t *testing.T) {
	// a
	// |- a1
	// |--- a11
	// |--- a12
	// |- a2
	// |--- a21
	// b
	// add
	// adddomain
	ret := []struct {
		Segments [][]rune
		Cands    [][]rune
		idx      int
		Ret      [][]rune
		pos      int
	}{
		{sr(""), sr("a", "b", "add", "adddomain"), 0, sr("a", "b", "add", "adddomain"), 0},
		{sr("a"), sr("a", "add", "adddomain"), 1, sr("", "dd", "dddomain"), 1},
		{sr("a", ""), sr("a1", "a2"), 0, sr("a1", "a2"), 0},
		{sr("a", "a"), sr("a1", "a2"), 1, sr("1", "2"), 1},
		{sr("a", "a1"), sr("a1"), 2, sr(""), 2},
		{sr("add"), sr("add", "adddomain"), 2, sr("", "domain"), 2},
	}
	for idx, r := range ret {
		cands, pos := RetSegment(r.Segments, r.Cands, r.idx)
		ret := make([][]rune, 0, len(cands))
		for _, v := range cands {
			ret = append(ret, v.Name)
		}

		if !reflect.DeepEqual(ret, r.Ret) {
			t.Errorf("incorrect ret %v at %d", ret, idx)
		}
		if pos != r.pos {
			t.Errorf("incorrect pos %d at %d", pos, idx)
		}
	}
}

func TestSplitSegment(t *testing.T) {
	// a
	// |- a1
	// |--- a11
	// |--- a12
	// |- a2
	// |--- a21
	// b
	ret := []struct {
		Line     string
		Pos      int
		Segments [][]rune
		Idx      int
	}{
		{"", 0, sr(""), 0},
		{"a", 1, sr("a"), 1},
		{"a ", 2, sr("a", ""), 0},
		{"a a", 3, sr("a", "a"), 1},
		{"a a1", 4, sr("a", "a1"), 2},
		{"a a1 ", 5, sr("a", "a1", ""), 0},
	}

	for i, r := range ret {
		ret, idx := SplitSegment([]rune(r.Line), r.Pos)
		if !reflect.DeepEqual(rs(ret), rs(r.Segments)) {
			t.Errorf("incorrect ret %v at %d", ret, i)
		}
		if idx != r.Idx {
			t.Errorf("incorrect idx %d at %d", idx, i)
		}
	}
}

type Tree struct {
	Name     string
	Children []Tree
}

func TestSegmentCompleter(t *testing.T) {
	tree := Tree{"", []Tree{
		{"a", []Tree{
			{"a1", []Tree{
				{"a11", nil},
				{"a12", nil},
			}},
			{"a2", []Tree{
				{"a21", nil},
			}},
		}},
		{"b", nil},
		{"route", []Tree{
			{"add", nil},
			{"adddomain", nil},
		}},
	}}
	s := SegmentFunc(func(ret [][]rune, n int) [][]rune {
		tree := tree
	main:
		for level := 0; level < len(ret)-1; {
			name := string(ret[level])
			for _, t := range tree.Children {
				if t.Name == name {
					tree = t
					level++
					continue main
				}
			}
		}

		ret = make([][]rune, len(tree.Children))
		for idx, r := range tree.Children {
			ret[idx] = []rune(r.Name)
		}
		return ret
	})

	// a
	// |- a1
	// |--- a11
	// |--- a12
	// |- a2
	// |--- a21
	// b
	ret := []struct {
		Line  string
		Pos   int
		Ret   [][]rune
		Share int
	}{
		{"", 0, sr("a", "b", "route"), 0},
		{"a", 1, sr(""), 1},
		{"a ", 2, sr("a1", "a2"), 0},
		{"a a", 3, sr("1", "2"), 1},
		{"a a1", 4, sr(""), 2},
		{"a a1 ", 5, sr("a11", "a12"), 0},
		{"a a1 a", 6, sr("11", "12"), 1},
		{"a a1 a1", 7, sr("1", "2"), 2},
		{"a a1 a11", 8, sr(""), 3},
		{"route add", 9, sr("", "domain"), 3},
	}
	for _, r := range ret {
		for idx, rr := range r.Ret {
			r.Ret[idx] = append(rr, ' ')
		}
	}
	for i, r := range ret {
		newLine, length := s.Do([]rune(r.Line), r.Pos, r.Pos)
		lines := make([][]rune, 0, len(newLine))
		for _, v := range newLine {
			lines = append(lines, v.Name)
		}

		if !reflect.DeepEqual(rs(lines), rs(r.Ret)) {
			t.Errorf("incorrect lines %v at %d", lines, i)
		}
		if length != r.Share {
			t.Errorf("incorrect length %d at %d", length, i)
		}
	}
}
