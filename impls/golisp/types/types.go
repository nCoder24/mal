package types

import (
	"fmt"
	"strings"
)

type MalValue interface{}

type Symbol string
type Int int
type String string
type Keyword string
type Func func(args []MalValue) (MalValue, error)
type List []MalValue
type Vector []MalValue
type Map []MalValue

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
