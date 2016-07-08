package pluralforms

import "fmt"

type Test interface {
	Test(n uint32) bool
	String() string
}

type Ternary struct {
	Test  Test
	True  Expression
	False Expression
}

func (t Ternary) Eval(n uint32) int {
	if t.Test.Test(n) {
		if t.True == nil {
			return -1
		}
		return t.True.Eval(n)
	} else {
		if t.False == nil {
			return -1
		}
		return t.False.Eval(n)
	}
}

func (t Ternary) String() string {
	return fmt.Sprintf("<Ternary(%s?%s:%s)>", t.Test, t.True, t.False)
}
