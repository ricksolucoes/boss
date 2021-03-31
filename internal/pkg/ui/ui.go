package ui

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

// GetTextOrDef get a new text or default value
func GetTextOrDef(label string, def string) string {
	prompt := promptui.Prompt{
		Label:       label,
		Default:     def,
		HideEntered: true,
		Pointer:     promptui.PipeCursor,
	}

	result, err := prompt.Run()

	if err != nil {
		if err.Error() == "^C" {
			os.Exit(1)
		}
		return def
	}

	return result
}

// GetText get a new text
func GetText(label string, required bool) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(result string) error {
			if len(result) == 0 && required {
				return fmt.Errorf("Value is required!")
			}
			return nil
		},
	}
	return prompt.Run()
}

// GetPassword prompts for credentials
func GetPassword(label string, required bool) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
		Validate: func(result string) error {
			if len(result) == 0 && required {
				return fmt.Errorf("Value is required!")
			}
			return nil
		},
	}
	return prompt.Run()
}

// GetConfirmation write a question
func GetConfirmation(question string, def bool) (bool, error) {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return def, nil
	}

	return result == "y", err
}
