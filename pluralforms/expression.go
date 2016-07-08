package pluralforms

import (
	"fmt"
	"strings"
)

type Expression interface {
	Eval(n uint32) int
	String() string
}

type Const struct {
	Value int
}

func (c Const) Eval(n uint32) int {
	return c.Value
}

func (c Const) String() string {
	return fmt.Sprintf("<Const:%d>", c.Value)
}

func pformat(expr Expression) string {
	ret := ""
	s := expr.String()
	level := -1
	for _, rune := range s {
		switch rune {
		case '<':
			level++
			ret += "\n" + strings.Repeat("  ", level)
		case '>':
			level--
			if level >= 0 {
				ret += "\n" + strings.Repeat("  ", level)
			}
		case ':', '?', '|':
		default:
			ret += fmt.Sprintf("%c", rune)
		}
	}
	return strings.Replace(ret, "\n\n", "\n", -1)
}
