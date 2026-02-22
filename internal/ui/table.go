package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// TableStyle defines the visual style for tables.
type TableStyle struct {
	HeaderStyle lipgloss.Style
	CellStyle   lipgloss.Style
	BorderStyle lipgloss.Style
}

// DefaultTableStyle returns the default styling for tables.
func DefaultTableStyle() TableStyle {
	return TableStyle{
		HeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1),
		CellStyle: lipgloss.NewStyle().
			Padding(0, 1),
		BorderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
	}
}

// Table wraps lipgloss table with minectl styling.
type Table struct {
	headers  []string
	rows     [][]string
	style    TableStyle
	headless bool
}

// NewTable creates a new styled table.
func NewTable(u *UI, headers ...string) *Table {
	return &Table{
		headers:  headers,
		rows:     make([][]string, 0),
		style:    DefaultTableStyle(),
		headless: u.IsHeadless(),
	}
}

// Append adds a row to the table.
func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

// buildTable creates the styled lipgloss table.
func (t *Table) buildTable() *table.Table {
	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(t.style.BorderStyle).
		Headers(t.headers...).
		Rows(t.rows...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return t.style.HeaderStyle
			}
			if row%2 == 0 {
				return t.style.CellStyle.Background(lipgloss.Color("235"))
			}
			return t.style.CellStyle
		})
}

// Render prints the table to stdout.
func (t *Table) Render() {
	if t.headless {
		t.renderPlain()
		return
	}
	fmt.Println(t.buildTable())
}

// renderPlain outputs a simple text table for headless mode.
func (t *Table) renderPlain() {
	fmt.Println(strings.Join(t.headers, "\t"))
	for _, row := range t.rows {
		fmt.Println(strings.Join(row, "\t"))
	}
}

// String returns the table as a string.
func (t *Table) String() string {
	if t.headless {
		var sb strings.Builder
		sb.WriteString(strings.Join(t.headers, "\t"))
		sb.WriteString("\n")
		for _, row := range t.rows {
			sb.WriteString(strings.Join(row, "\t"))
			sb.WriteString("\n")
		}
		return sb.String()
	}
	return t.buildTable().String()
}

// RenderToWriter writes the table to the given writer.
func (t *Table) RenderToWriter(w io.Writer) {
	_, _ = fmt.Fprintln(w, t.String())
}
