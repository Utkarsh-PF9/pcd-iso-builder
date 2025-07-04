package main

import tea "github.com/charmbracelet/bubbletea"

type storage struct{}

func (n storage) Init() tea.Cmd                           { return nil }
func (n storage) Update(msg tea.Msg) (tea.Model, tea.Cmd) { 
	switch msg:=msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case "esc":
			return n, func() tea.Msg {
				return switchPageMsg("back")
			} 
		}
	}
	return n, nil 
}
func (n storage) View() string                            { return "Storage Configuration Page" }

func StorageInitialModel() storage {
	return storage{}
}
