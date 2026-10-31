package ui

import "testing"

func TestFormTheme(t *testing.T) {
	tests := []struct {
		name   string
		isDark bool
	}{
		{"dark background", true},
		{"light background", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			styles := FormTheme().Theme(tt.isDark)
			if styles == nil {
				t.Fatal("FormTheme returned nil styles")
			}
			if !styles.Focused.Title.GetBold() {
				t.Error("expected focused title to be bold")
			}
		})
	}
}
