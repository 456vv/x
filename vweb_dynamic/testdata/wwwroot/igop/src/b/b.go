package b

import (
	"fmt"
)

func Sprint(m any) string {
	return fmt.Sprint(m, "+receive")
}
