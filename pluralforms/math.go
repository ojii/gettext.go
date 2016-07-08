package pluralforms

import "fmt"

type Math interface {
	Calc(n uint32) uint32
	String() string
}

type Mod struct {
	Value uint32
}

func (m Mod) Calc(n uint32) uint32 {
	return n % m.Value
}

func (m Mod) String() string {
	return fmt.Sprintf("<Mod(%d)>", m.Value)
}
