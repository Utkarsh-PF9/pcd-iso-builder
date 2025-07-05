package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var stack = []string{"menu"}

var helpMap map[string]string = map[string]string{
	"menu":    "move: ↑/↓  •  select: enter  •  exit: ctrl + c",
	"network": "next field: tab  •  previous field: shift + tab  •  submit: enter  •  back: esc  •  exit: ctrl + c",
	"storage": "back: esc  •  exit: ctrl + c",
	"form":    "next field: tab  •  previous field: shift + tab  •  submit: enter  •  back: esc  •  exit: ctrl + c",
}

type layout struct {
	width, height int
	currentPage   tea.Model
}

var isConfigured bool = false

func checkIsPCDCtlConfigured() bool {
	return isConfigured
}

func setIsPCDCtlConfigured(val bool) {
	isConfigured = val
}

func (m layout) Init() tea.Cmd {

	// helpMap["menu"]="move: ↑/↓  •  select: enter  •  exit: ctrl + c"
	// helpMap["network"]="back: esc  •  exit: ctrl + c"
	// helpMap["storage"]="back: esc  •  exit: ctrl + c"

	return m.currentPage.Init()
}

func (m layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	backFunction := func() {
		if len(stack) > 1 {
			stack = stack[:len(stack)-1]
		}
		page := stack[len(stack)-1]

		switch page {
		case "menu":
			m.currentPage = MenuInitialModel(checkIsPCDCtlConfigured())
		case "network":
			m.currentPage = NetworkInitialModel(m.width, m.height)
		case "storage":
			m.currentPage = StorageInitialModel()
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = int(0.7 * float64(msg.Width))
		m.height = int(0.7 * float64(msg.Height))

		// Propagate new size to current page if it supports WindowSizeMsg
		if wm, ok := m.currentPage.(tea.Model); ok {
			m.currentPage, _ = wm.Update(msg)
		}

	case configurationStatus:
		switch msg {
		case "done":
			setIsPCDCtlConfigured(true)
			backFunction()
		}

	case switchPageMsg:
		switch msg {
		case "form":
			stack = append(stack, "form")
			m.currentPage = FormInitialModel(m.width,m.height)
		case "network":
			stack = append(stack, "network")
			m.currentPage = NetworkInitialModel(m.width, m.height)
		case "storage":
			stack = append(stack, "storage")
			m.currentPage = StorageInitialModel()

		case "back":
			backFunction()
		}
		return m, m.currentPage.Init()

	}

	var cmd tea.Cmd
	m.currentPage, cmd = m.currentPage.Update(msg)
	return m, cmd
}

func (m layout) View() string {

	headerList := ""
	for _, v := range stack {
		headerList = headerList + " › " + v
	}

	header := lipgloss.NewStyle().
		Align(lipgloss.Left, lipgloss.Center).
		Width(m.width).
		Padding(1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Foreground(lipgloss.Color("15")).
		Render(headerList)

	// footerHelp:=helpMap["menu"]

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Width(m.width).
		Padding(1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		Foreground(lipgloss.Color("15")).
		Render(helpMap[stack[len(stack)-1]])

	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
	contentWidth := m.width

	h, v := docStyle.GetFrameSize()

	// If currentPage is a list.Model, update its size
	if listPage, ok := m.currentPage.(menu); ok {
		listPage.list.SetSize(contentWidth-h, contentHeight-v)
		m.currentPage = listPage
	}

	content := lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).
		Width(contentWidth).
		Height(contentHeight).
		Foreground(lipgloss.Color("15")).
		Render(m.currentPage.View())

	layoutBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("15")).
		// Padding(1).
		Render(lipgloss.JoinVertical(lipgloss.Top, header, content, footer))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		layoutBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

func LayoutInitialModel(width, height int) layout {
	return layout{
		width:       int(0.7 * float64(width)),
		height:      int(0.7 * float64(height)),
		currentPage: MenuInitialModel(checkIsPCDCtlConfigured()),
	}
}
