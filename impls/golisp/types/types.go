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

func stringify(forms []MalValue) string {
	strs := make([]string, 0, len(forms))

	for _, v := range forms {
		strs = append(strs, fmt.Sprintf("%v", v))
	}

	return strings.Join(strs, " ")
}

func Seq(mal MalValue) ([]MalValue, bool) {
	switch v := mal.(type) {
	case List:
		return v, true
	case Vector:
		return v, true
	case Map:
		return v, true
	}

	return nil, false
}
