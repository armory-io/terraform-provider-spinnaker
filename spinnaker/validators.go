package spinnaker

import (
	"fmt"
	"regexp"
)

func validateApplicationName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Only alphanumeric characters or '-' allowed in %q", k))
	}
	return
}
