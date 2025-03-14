package utils

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

const Indentation = "  "

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func TrimAndIndent(s string) string {
	s = Trim(s)
	lines := strings.Split(s, "\n")
	indentedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmed := Trim(line)
		indentedLines[i] = Indentation + trimmed
	}
	return strings.Join(indentedLines, "\n")
}

func SafeString(s *string) string {
	if s == nil {
		return "" // or "<nil>" or whatever default you want
	}
	return *s
}

func SafeInt(i *int) int {
	if i == nil {
		return 0 // or some other default
	}
	return *i
}

func GetApproval() (bool, error) {
	templates := &promptui.SelectTemplates{
		// Help: "Use the arrow keys to navigate: ↓ ↑ → ←",
		// Label: "{{ . }}",
		// Active:   fmt.Sprintf(`{{ "Do you approve" | green }} {{ "%v:" | faint }} {{ . | faint }}`, promptui.IconGood),
		// Inactive: fmt.Sprintf(`{{ "Do you approve" | green }} {{ "%v:" | faint }} {{ . | faint }}`, promptui.IconGood),
		Selected: fmt.Sprintf(`{{ "Do you approve" | green }} {{ "%v:" | faint }} {{ . | faint }}`, promptui.IconGood),
	}

	prompt := promptui.Select{
		Label:     "Do you approve",
		Items:     []string{"No", "Yes"},
		Templates: templates,
		// Stdout:    NoBellStdout,
	}

	_, response, err := prompt.Run()
	if err != nil {
		return false, err
	}

	approved := response == "Yes"
	return approved, nil
}
