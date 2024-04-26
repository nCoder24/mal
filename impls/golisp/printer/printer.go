package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nCoder24/mal/impls/golisp/types"
)

func PrStr(mal types.MalValue, quoted bool) string {
	switch v := mal.(type) {
	case types.String:
		if quoted {
			return strconv.Quote(string(v))
		}

		return string(v)
	case types.List:
		return "(" + listStr(v, quoted) + ")"
	case types.Vector:
		return "[" + listStr(v, quoted) + "]"
	case types.Map:
		return "{" + listStr(v, quoted) + "}"
	}

	return fmt.Sprintf("%v", mal)
}

func listStr(mals []types.MalValue, quoted bool) string {
	strs := make([]string, 0, len(mals))

	for _, mal := range mals {
		strs = append(strs, PrStr(mal, quoted))
	}

	return strings.Join(strs, " ")
}
