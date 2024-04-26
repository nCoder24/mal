package core

import (
	"github.com/nCoder24/mal/impls/golisp/types"
)

var Namespace = map[string]types.MalValue{
	"+":       types.Func(sum),
	"-":       types.Func(sub),
	"*":       types.Func(mul),
	"/":       types.Func(divide),
	"=":       types.Func(equal),
	"<":       types.Func(lt),
	">":       types.Func(gt),
	"<=":      types.Func(lte),
	">=":      types.Func(gte),
	"list":    types.Func(list),
	"list?":   types.Func(isList),
	"empty?":  types.Func(isEmpty),
	"count":   types.Func(count),
	"prn":     types.Func(prn),
	"pr-str":  types.Func(prStr),
	"str":     types.Func(str),
	"println": types.Func(printLine),
}
