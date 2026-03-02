package handlers

import (
	"fmt"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// Return string for describing error and unique status code
func buildErrorsDocumentation(errors []huma.StatusError) (errDoc string, errCodes []int) {
	var errDocs []string
	statusMap := make(map[int]bool)

	for _, e := range errors {
		status := e.GetStatus()
		errDocs = append(errDocs, fmt.Sprintf("- `%d`: %s", status, e.Error()))

		if !statusMap[status] {
			statusMap[status] = true
			errCodes = append(errCodes, status)
		}
	}

	errDoc = fmt.Sprintf("\n### Error messages\n%s", strings.Join(errDocs, "\n"))

	return errDoc, errCodes
}
