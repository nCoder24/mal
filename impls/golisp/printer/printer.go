package printer

import (
	"fmt"

	"github.com/nCoder24/mal/impls/golisp/types"
)

func PrStr(mal types.MalValue) string {
	if mal == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", mal)
}
