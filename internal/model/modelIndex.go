package model

import "regexp"

func IsIntegerVariable(variable string) bool {
	// Check if the variable is not empty
	if variable == "" {
		return false
	}

	// Use a regular expression to match only integer values
	match, _ := regexp.MatchString(`^\d+$`, variable)

	return match
}
