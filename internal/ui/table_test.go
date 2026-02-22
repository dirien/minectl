package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableAppend(t *testing.T) {
	u := newTestUI(t, true)
	tbl := NewTable(u, "ID", "NAME")
	tbl.Append([]string{"1", "test-server"})
	tbl.Append([]string{"2", "other-server"})

	if got := len(tbl.rows); got != 2 {
		t.Errorf("got %d rows, want 2", got)
	}
}

func TestTableString(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
		headers  []string
		rows     [][]string
		wantAll  []string
	}{
		{
			name:     "headless tab-separated output",
			headless: true,
			headers:  []string{"ID", "NAME", "REGION"},
			rows: [][]string{
				{"abc-123", "my-server", "fra1"},
				{"def-456", "other-server", "lon1"},
			},
			wantAll: []string{
				"ID\tNAME\tREGION",
				"abc-123\tmy-server\tfra1",
				"def-456\tother-server\tlon1",
			},
		},
		{
			name:     "interactive styled output",
			headless: false,
			headers:  []string{"ID", "NAME"},
			rows:     [][]string{{"1", "test"}},
			wantAll:  []string{"ID", "NAME", "test"},
		},
		{
			name:     "empty table with headers only",
			headless: true,
			headers:  []string{"ID", "NAME"},
			rows:     nil,
			wantAll:  []string{"ID\tNAME"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			tbl := NewTable(u, tt.headers...)
			for _, row := range tt.rows {
				tbl.Append(row)
			}

			result := tbl.String()
			for _, want := range tt.wantAll {
				if !strings.Contains(result, want) {
					t.Errorf("output missing %q\ngot: %q", want, result)
				}
			}
		})
	}
}

func TestTableRenderToWriter(t *testing.T) {
	u := newTestUI(t, true)
	tbl := NewTable(u, "ID", "NAME")
	tbl.Append([]string{"1", "test"})

	var buf bytes.Buffer
	tbl.RenderToWriter(&buf)

	if !strings.Contains(buf.String(), "1\ttest") {
		t.Errorf("RenderToWriter output missing expected content, got %q", buf.String())
	}
}

func TestDefaultTableStyle(t *testing.T) {
	style := DefaultTableStyle()
	if !style.HeaderStyle.GetBold() {
		t.Error("expected header style to be bold")
	}
}
