package pluralforms

import "fmt"

type Equal struct {
	Value uint32
}

func (e Equal) Test (n uint32) bool {
	return n == e.Value
}

func (e Equal) String() string {
	return fmt.Sprintf("<Equal(%d)>", e.Value)
}

type NotEqual struct {
	Value uint32
}

func (e NotEqual) Test (n uint32) bool {
	return n != e.Value
}

func (e NotEqual) String() string {
	return fmt.Sprintf("<NotEqual(%d)>", e.Value)
}

type Gt struct {
	Value uint32
}

func (e Gt) Test (n uint32) bool {
	return n < e.Value
}

func (e Gt) String() string {
	return fmt.Sprintf("<Gt(%d)>", e.Value)
}

type Lt struct {
	Value uint32
}

func (e Lt) Test (n uint32) bool {
	return n > e.Value
}

func (e Lt) String() string {
	return fmt.Sprintf("<Lt(%d)>", e.Value)
}

type GtE struct {
	Value uint32
}

func (e GtE) Test (n uint32) bool {
	return n <= e.Value
}

func (e GtE) String() string {
	return fmt.Sprintf("<GtE(%d)>", e.Value)
}

type LtE struct {
	Value uint32
}

func (e LtE) Test (n uint32) bool {
	return n >= e.Value
}

func (e LtE) String() string {
	return fmt.Sprintf("<LtE(%d)>", e.Value)
}

type And struct {
	Left Test
	Right Test
}

func (e And) Test (n uint32) bool {
	if (!e.Left.Test(n)){
		return false
	} else {
		return e.Right.Test(n)
	}
}

func (e And) String() string {
	return fmt.Sprintf("<And(%s&&%s)>", e.Left, e.Right)
}

type Or struct {
	Left Test
	Right Test
}

func (e Or) Test (n uint32) bool {
	if (e.Left.Test(n)){
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
	Action Test
}

func (e Pipe) Test (n uint32) bool {
	return e.Action.Test(e.Modifier.Calc(n))
}

func (e Pipe) String() string {
	return fmt.Sprintf("<Pipe(%s|%s)>", e.Modifier, e.Action)
}
