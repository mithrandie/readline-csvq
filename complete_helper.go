package readline

import (
	"bytes"
	"strings"
	"unicode"
)

// Caller type for dynamic completion
type DynamicCompleteFunc func(string, string, int) ([]string, bool)

type PrefixCompleterInterface interface {
	Print(prefix string, level int, buf *bytes.Buffer)
	Do(line []rune, pos int, index int) (newLine [][]rune, length int, formatAsIdentifier bool)
	GetName() ([]rune, bool)
	GetChildren() []PrefixCompleterInterface
	SetChildren(children []PrefixCompleterInterface)
	IsAppendOnly() bool
}

type DynamicPrefixCompleterInterface interface {
	PrefixCompleterInterface
	IsDynamic() bool
	GetDynamicNames(line []rune, origLine []rune, index int) ([][]rune, bool)
}

type PrefixCompleter struct {
	Name               []rune
	Dynamic            bool
	Callback           DynamicCompleteFunc
	Children           []PrefixCompleterInterface
	FormatAsIdentifier bool
	AppendOnly         bool
}

func (p *PrefixCompleter) Tree(prefix string) string {
	buf := bytes.NewBuffer(nil)
	p.Print(prefix, 0, buf)
	return buf.String()
}

func Print(p PrefixCompleterInterface, prefix string, level int, buf *bytes.Buffer) {
	name, _ := p.GetName()
	if strings.TrimSpace(string(name)) != "" {
		buf.WriteString(prefix)
		if level > 0 {
			buf.WriteString("├")
			buf.WriteString(strings.Repeat("─", (level*4)-2))
			buf.WriteString(" ")
		}
		buf.WriteString(string(name) + "\n")
		level++
	}
	for _, ch := range p.GetChildren() {
		ch.Print(prefix, level, buf)
	}
}

func (p *PrefixCompleter) Print(prefix string, level int, buf *bytes.Buffer) {
	Print(p, prefix, level, buf)
}

func (p *PrefixCompleter) IsDynamic() bool {
	return p.Dynamic
}

func (p *PrefixCompleter) GetName() ([]rune, bool) {
	return p.Name, p.FormatAsIdentifier
}

func (p *PrefixCompleter) GetDynamicNames(line []rune, origLine []rune, index int) ([][]rune, bool) {
	c, formatAsIdentifier := p.Callback(string(line), string(origLine), index)
	names := make([][]rune, 0, len(c))
	for _, name := range c {
		names = append(names, []rune(name+" "))
	}
	return names, formatAsIdentifier
}

func (p *PrefixCompleter) GetChildren() []PrefixCompleterInterface {
	return p.Children
}

func (p *PrefixCompleter) SetChildren(children []PrefixCompleterInterface) {
	p.Children = children
}

func (p *PrefixCompleter) IsAppendOnly() bool {
	return p.AppendOnly
}

func NewPrefixCompleter(pc ...PrefixCompleterInterface) *PrefixCompleter {
	return PcItem("", pc...)
}

func PcItem(name string, pc ...PrefixCompleterInterface) *PrefixCompleter {
	name += " "
	return &PrefixCompleter{
		Name:               []rune(name),
		Dynamic:            false,
		Children:           pc,
		FormatAsIdentifier: false,
		AppendOnly:         false,
	}
}

func PcItemDynamic(callback DynamicCompleteFunc, pc ...PrefixCompleterInterface) *PrefixCompleter {
	return &PrefixCompleter{
		Callback:           callback,
		Dynamic:            true,
		Children:           pc,
		FormatAsIdentifier: false,
		AppendOnly:         false,
	}
}

func (p *PrefixCompleter) Do(line []rune, pos int, index int) (newLine [][]rune, offset int, formatAsIdentifier bool) {
	return doInternal(p, line, pos, line, index)
}

func Do(p PrefixCompleterInterface, line []rune, pos int, index int) (newLine [][]rune, offset int, formatAsIdentifier bool) {
	return doInternal(p, line, pos, line, index)
}

func doInternal(p PrefixCompleterInterface, line []rune, pos int, origLine []rune, index int) (newLine [][]rune, offset int, formatAsIdentifier bool) {
	line = runes.TrimSpaceLeft(line[:pos])
	goNext := false
	var lineCompleter PrefixCompleterInterface
	for _, child := range p.GetChildren() {
		childNames := make([][]rune, 1)

		if child.IsAppendOnly() {
			line = []rune(LastElement(string(line)))
		}

		childDynamic, ok := child.(DynamicPrefixCompleterInterface)
		if ok && childDynamic.IsDynamic() {
			childNames, formatAsIdentifier = childDynamic.GetDynamicNames(line, origLine, index)
		} else {
			childNames[0], formatAsIdentifier = child.GetName()
		}

		for _, childName := range childNames {
			if len(line) >= len(childName) {
				if runes.HasPrefixFold(line, childName, formatAsIdentifier) {
					if len(line) == len(childName) {
						newLine = append(newLine, append(childName, ' '))
					} else {
						newLine = append(newLine, childName)
					}
					offset = len(childName)
					lineCompleter = child
					goNext = true
				}
			} else {
				if runes.HasPrefixFold(childName, line, formatAsIdentifier) {
					newLine = append(newLine, childName)
					offset = len(line)
					lineCompleter = child
				}
			}
		}
	}

	if len(newLine) != 1 {
		return
	}

	tmpLine := make([]rune, 0, len(line))
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}

		tmpLine = append(tmpLine, line[i:]...)
		return doInternal(lineCompleter, tmpLine, len(tmpLine), origLine, index)
	}

	if goNext {
		return doInternal(lineCompleter, nil, 0, origLine, index)
	}
	return
}

func LastElement(input string) string {
	src := []rune(input)
	buf := new(bytes.Buffer)

	var quote rune = -1

	for i := 0; i < len(src); i++ {
		if 0 < buf.Len() || -1 < quote {
			if -1 < quote && src[i] == '\\' {
				if i+1 < len(src) && src[i+1] == quote {
					i++
				}
			}

			if (-1 < quote && src[i] == quote) || (quote == -1 && unicode.IsSpace(src[i])) {
				quote = -1
				buf.Reset()
			} else {
				buf.WriteRune(src[i])
			}
			continue
		}

		if unicode.IsSpace(src[i]) {
			continue
		}

		if src[i] == '\'' || src[i] == '"' || src[i] == '`' {
			quote = src[i]
		}

		buf.WriteRune(src[i])
	}
	return buf.String()
}
