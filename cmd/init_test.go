package cmd

import (
	"testing"
)

func TestInit_InvalidName(t *testing.T) {
	cases := []string{"", ".", "..", "../..", "/"}
	initCmd := newInitCmd()

	for _, name := range cases {
		err := initCmd.ExecuteArgs(t.Context(), []string{name})
		if err == nil {
			t.Errorf("expected error for name %s, got nil", name)
		}
	}
}