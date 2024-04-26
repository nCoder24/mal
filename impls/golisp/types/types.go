package types

import (
	"go/types"
	"strconv"
)

var (
	Nil   = NilPtr{}
	True  = Bool(true)
	False = Bool(false)
)

type MalValue interface {
}

type NilPtr types.Nil
type String string
type Bool bool
type Symbol string
type Number float64
type Keyword string
type List []MalValue
type Vector []MalValue
type Map []MalValue
type Func func(args []MalValue) (MalValue, error)

type Env interface {
}

type DefinedFunc struct {
	Func
	Env      any
	Bindings []string
	Body     MalValue
}

func (n NilPtr) String() string {
	return "nil"
}

func (f Func) String() string {
	return "#<function>"
}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}
