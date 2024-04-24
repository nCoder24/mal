package types

import (
	"fmt"
	"go/types"
	"strings"
)

var (
	Nil   = NilPtr{}
	True  = Bool(true)
	False = Bool(false)
)

type MalValue interface{}

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

func stringify(forms []MalValue) string {
	strs := make([]string, 0, len(forms))

	for _, v := range forms {
		strs = append(strs, fmt.Sprintf("%v", v))
	}

	return strings.Join(strs, " ")
}

func Seq(val MalValue) ([]MalValue, error) {
	switch v := val.(type) {
	case List:
		return v, nil
	case Vector:
		return v, nil
	case Map:
		return v, nil
	}

	return nil, fmt.Errorf("cannot convert '%v' to seq", val)
}
