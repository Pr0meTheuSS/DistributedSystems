package handler

import (
	"fmt"
)

func buildErrorMessage(err error, code int) string {
	return fmt.Sprintf("Error: %s. Status code: %d", err.Error(), code)
}
