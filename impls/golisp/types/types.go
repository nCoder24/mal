package types

import (
	"go/types"
)

var (
	Nil   = NilPtr{}
	True  = Bool(true)
	False = Bool(false)
)

type MalValue interface {
}

type Bool bool
type NilPtr types.Nil
type Symbol string
type Number float64
type String string
type Keyword string
type Func func(args []MalValue) (MalValue, error)
type List []MalValue
type Vector []MalValue
type Map []MalValue

func (n NilPtr) String() string {
	return "nil"
}

func (l List) String() string {
	return "(" + stringify(l) + ")"
}

func (v Vector) String() string {
	return "[" + stringify(v) + "]"
}

func (v Map) String() string {
	return "{" + stringify(v) + "}"
}

func (f Func) String() string {
	return "#<function>"
}
