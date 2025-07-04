package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type outer struct {
	width, height int
	child         tea.Model // layout or loading screen
}



func (o outer) Init() tea.Cmd {
	return o.child.Init()
}



func (o outer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Resize message, capture terminal size
	switch msg := msg.(type) {
	case switchPageMsg:
		switch msg {
		case "menu":
			o.child = LayoutInitialModel(o.width, o.height)
			return o, o.child.Init()
		}
	case tea.WindowSizeMsg:
		o.width = msg.Width
		o.height = msg.Height

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return o, tea.Quit
		}
	}

	var cmd tea.Cmd
	o.child, cmd = o.child.Update(msg)
	return o, cmd
}



func (o outer) View() string {
	return lipgloss.NewStyle().
		Width(o.width).
		Height(o.height).
		Background(lipgloss.Color("0")).
		Render(
			lipgloss.Place(
				o.width, o.height,
				lipgloss.Center, lipgloss.Center,
				o.child.View(),
			),
		)
}
