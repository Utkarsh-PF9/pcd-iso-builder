package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadingModel struct {
	width, height int
	showCursor    bool
}

func (m loadingModel) Init() tea.Cmd {
	return tea.Batch(blinkCursor(), transitionToOuter())
}

type blinkMsg struct{}
type transitionMsg struct{}



func blinkCursor() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return blinkMsg{}
	})
}



func transitionToOuter() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return transitionMsg{}
	})
}



func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case blinkMsg:
		m.showCursor = !m.showCursor
		return m, blinkCursor()

	case transitionMsg:

		return outer{
			width:  m.width,
			height: m.height,
			child: LayoutInitialModel(m.width, m.height),  // CHECK WHETHER THE PCDCTL IS ALREADY CONFIGURED AND THEN REDIRECT ACCORDINGLY
			// child: FormInitialModel(),
		}, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}



func (m loadingModel) View() string {
	// "Platform" styled in bold white
	platformStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15"))

	// "9" styled in bold blue with blinking background
	nineStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")) // blue foreground

	if m.showCursor {
		nineStyle = nineStyle.Background(lipgloss.Color("39")) // white bg
	} else {
		nineStyle = nineStyle.Background(lipgloss.Color("0")) // black bg
	}

	// Combine both parts
	text := platformStyle.Render("Platform") + nineStyle.Render("9")

	// Render centered
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(lipgloss.Color("0")). // overall black background
		Render(
			lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				text,
			),
		)
}
