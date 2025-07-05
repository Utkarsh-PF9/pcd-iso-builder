package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var items = []list.Item{
	item{title: "Configure Network", desc: "IP, Gateway, DNS, Subnet, NAT, Bonds"},
	item{title: "Configure Storage", desc: "Resize, Partition, etc"},
}

var setupItems = []list.Item{
	item{title: "Configure Network", desc: "IP, Gateway, DNS, Subnet, NAT"},
	item{title: "Configure PCDCtl", desc: "Installation and configuration of required packages"},
}

type switchPageMsg string
type configurationStatus string

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type menu struct {
	isConfigured bool
	list         list.Model
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.isConfigured {
				switch m.list.Index() {
				case 0:
					return m, func() tea.Msg {
						return switchPageMsg("network")
					}
				case 1:
					return m, func() tea.Msg {
						return switchPageMsg("storage")
					}
				}
			} else {
				switch m.list.Index() {
				case 0:
					return m, func() tea.Msg {
						return switchPageMsg("network")
					}
				case 1:
					return m, func() tea.Msg {
						return switchPageMsg("form")
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg) // Let list handle all key events
	return m, cmd
}

func (m menu) View() string {
	return docStyle.Render(m.list.View())
}

func MenuInitialModel(isConfigured bool) menu {
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("39")).BorderForeground(lipgloss.Color("39")) // Blue for selected item
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("39")).BorderForeground(lipgloss.Color("39"))

	m := menu{list: list.New(items, delegate, 0, 0),isConfigured: isConfigured}

	if !isConfigured {
		m = menu{list: list.New(setupItems, delegate, 0, 0),isConfigured: isConfigured}
	}
	m.list.SetShowTitle(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.SetShowHelp(false)
	m.list.KeyMap = list.KeyMap{
		CursorUp:   key.NewBinding(key.WithKeys("up")),
		CursorDown: key.NewBinding(key.WithKeys("down")),
	}

	return m
}
