package commands

import (
	"testing"
)

func TestCommandFactory(t *testing.T) {
	debugFlag := true
	app := NewApplication([]string{"program", "command"}, func(int) {})

	testCases := []struct {
		name       string
		commandKey string
		expected   string
	}{
		{
			name:       "TestLinter",
			commandKey: "lint",
			expected:   "Linter",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmdFactory := app.commands[tc.commandKey]
			cmd, err := cmdFactory(&debugFlag)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cmd == nil {
				t.Errorf("CommandFactory returned nil command for key: %s", tc.commandKey)
			}
			if cmd.Name() != tc.expected {
				t.Errorf("Expected command name %s, but got %s", tc.expected, cmd.Name())
			}
		})
	}
}
