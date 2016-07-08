package pluralforms

import "fmt"

type Equal struct {
	Value uint32
}

func (e Equal) Test(n uint32) bool {
	return n == e.Value
}

func (e Equal) String() string {
	return fmt.Sprintf("<Equal(%d)>", e.Value)
}

type NotEqual struct {
	Value uint32
}

func (e NotEqual) Test(n uint32) bool {
	return n != e.Value
}

func (e NotEqual) String() string {
	return fmt.Sprintf("<NotEqual(%d)>", e.Value)
}

type Gt struct {
	Value   uint32
	Flipped bool
}

func (e Gt) Test(n uint32) bool {
	if e.Flipped {
		return e.Value > n
	} else {
		return n > e.Value
	}
}

func (e Gt) String() string {
	return fmt.Sprintf("<Gt(%d,%t)>", e.Value, e.Flipped)
}

type Lt struct {
	Value   uint32
	Flipped bool
}

func (e Lt) Test(n uint32) bool {
	if e.Flipped {
		return e.Value < n
	} else {
		return n < e.Value
	}
}

func (e Lt) String() string {
	return fmt.Sprintf("<Lt(%d,%t)>", e.Value, e.Flipped)
}

type GtE struct {
	Value   uint32
	Flipped bool
}

func (e GtE) Test(n uint32) bool {
	if e.Flipped {
		return e.Value >= n
	} else {
		return n >= e.Value
	}
}

func (e GtE) String() string {
	return fmt.Sprintf("<GtE(%d,%t)>", e.Value, e.Flipped)
}

type LtE struct {
	Value   uint32
	Flipped bool
}

func (e LtE) Test(n uint32) bool {
	if e.Flipped {
		return e.Value <= n
	} else {
		return n <= e.Value
	}
}

func (e LtE) String() string {
	return fmt.Sprintf("<LtE(%d,%t)>", e.Value, e.Flipped)
}

type And struct {
	Left  Test
	Right Test
}

func (e And) Test(n uint32) bool {
	if !e.Left.Test(n) {
		return false
	} else {
		return e.Right.Test(n)
	}
}

func (e And) String() string {
	return fmt.Sprintf("<And(%s&&%s)>", e.Left, e.Right)
}

type Or struct {
	Left  Test
	Right Test
}

func (e Or) Test(n uint32) bool {
	if e.Left.Test(n) {
		return true
	} else {
		return e.Right.Test(n)
	}
}

func (e Or) String() string {
	return fmt.Sprintf("<Or(%s||%s)>", e.Left, e.Right)
}

type Pipe struct {
	Modifier Math
	Action   Test
}

func (e Pipe) Test(n uint32) bool {
	return e.Action.Test(e.Modifier.Calc(n))
}

func (e Pipe) String() string {
	return fmt.Sprintf("<Pipe(%s|%s)>", e.Modifier, e.Action)
}
