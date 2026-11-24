package database

import (
	"fmt"
	"slices"
)

// string functions

// the msg is an optional field that tells the user whats wrong with the input
type validateStringInput func(in string) (correct bool, msg string)

var currentStringValidator validateStringInput = defaultStringValidator

// accept any argument
var defaultStringValidator validateStringInput = func(in string) (correct bool, msg string) {
	return true, "✓ Valid"

	//return false, fmt.Sprintf("✗ '%s' is not in the allowed list.", trimmed)
}

// this function overwrites the currentStringValidator so the user can define their own validators
// this action can be reseted using the ResetValidateStringInput function
func SetValidateStringInput(f validateStringInput) {
	currentStringValidator = f
}

// this function overwrites the currentStringValidator to the default validator
func ResetValidateStringInput(f validateStringInput) {
	currentStringValidator = f
}

func ValidateFloat(in string) (correct bool, msg string) {
	for _, c := range in {
		if !((c >= '0' && c <= '9') || c == '-' || c == '_' || c == ',' || c == '.') {
			return false, fmt.Sprintf("✗ '%c' is not a valid character.", c)
		}
	}
	return true, "✓ Valid Number"
}

// TODO: add documentation
func ValidateBasedOnOptions(options []string) {
	currentStringValidator = func(in string) (correct bool, msg string) {
		if slices.Contains(options, in) {
			return true, "✓ Valid"
		}
		return false, fmt.Sprintf("✗ '%s' is not in the allowed list.", in)
	}
}
