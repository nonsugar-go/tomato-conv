package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nonsugar-go/tools/tui"
)

// DevTypeList is a selection list of DevTypes.
func (c ConfInfo) DevTypeList() DevType {
	items := []list.Item{
		tui.Item("FortiGate"),
		tui.Item("PaloAlto"),
	}

	const defaultWidth = 20

	l := list.New(items, tui.ItemDelegate{}, defaultWidth, tui.ListHeight)
	l.Title = "機器の種類を選択してください"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = tui.TitleStyle
	l.Styles.PaginationStyle = tui.PaginationStyle
	l.Styles.HelpStyle = tui.HelpStyle

	m := tui.Model{List: l}

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	switch result.(tui.Model).List.SelectedItem().(tui.Item) {
	case "FortiGate":
		return DevTypeFortiGate
	case "PaloAlto":
		return DevTypePaloAlto
	}
	return DevTypeUnknown
}
