package b

import (
	"fmt"
)

func Sprint(m string) string {
	return fmt.Sprint(m, "+receive")
}