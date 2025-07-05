package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type storage struct {
	content string
}

func (n storage) Init() tea.Cmd {
	// Load the file contents at startup
	return func() tea.Msg {
		data, err := os.ReadFile("/etc/netplan/pf9_netplan.yaml")
		if err != nil {
			return "Storage Configuration Page"
		}
		return string(data)
	}
}

func (n storage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		n.content = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return n, func() tea.Msg {
				return switchPageMsg("back")
			}
		}
	}
	return n, nil
}

func (n storage) View() string {
	// Ensure it’s nicely displayed even for long content
	if strings.TrimSpace(n.content) == "" {
		return "Storage Configuration Page"
	}
	return n.content
}

func StorageInitialModel() storage {
	return storage{content: "Loading..."}
}
