package cli

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptString asks the user for a single string value.
func PromptString(label, defaultValue string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

// PromptSelect asks the user to select an option from a list.
func PromptSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

// ValidateNotEmpty is an optional validator that ensures the input is not entirely whitespace.
func ValidateNotEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("input cannot be empty")
	}
	return nil
}
