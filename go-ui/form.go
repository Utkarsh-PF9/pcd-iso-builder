package main

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type screenState int

const (
	fillForm screenState = iota
)

type setup struct {
	state    screenState
	form     *huh.Form
	formData values
}

type values struct {
	url      string
	username string
	region   string
	tenant   string
	password string
	confirm  bool
}

func (s setup) Init() tea.Cmd {
	return s.form.Init()
}

func (s setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch s.state {
	case fillForm:
		form, cmd := s.form.Update(msg)
		if updated, ok := form.(*huh.Form); ok {
			s.form = updated
		}

		if s.form.GetBool("CONFIRM") && s.form.State == huh.StateCompleted {
			return s, func() tea.Msg { return configurationStatus("done") }
		}
		return s, cmd
	}
	return s, nil
}

func (s setup) View() string {
	return s.form.View()
}

func FormInitialModel() setup {
	t := huh.ThemeCharm()

	t.Focused.Base = t.Focused.Base.Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("33"))
	t.Blurred.Base = t.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, true, false)
	t.Blurred.Title = t.Blurred.Title.Foreground(lipgloss.Color("#3c3c3c")).Padding(0, 0, 1, 0)
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("33")).Padding(0, 0, 1, 0)
	t.Group.Title = t.Group.Title.Foreground(lipgloss.Color("#ffffff")).Bold(true).Padding(0, 0, 2, 0)
	t.Form.Base = t.Form.Base.Padding(3).AlignHorizontal(lipgloss.Center)

	buttonTheme := huh.ThemeCharm()

	buttonTheme.Blurred.Title = buttonTheme.Blurred.Title.Foreground(lipgloss.Color("33"))
	buttonTheme.Focused.Title = buttonTheme.Focused.Title.Foreground(lipgloss.Color("33"))

	buttonTheme.Blurred.Title = buttonTheme.Blurred.Title.Margin(0, 0, 1, 0)
	buttonTheme.Focused.Title = buttonTheme.Focused.Title.Margin(0, 0, 1, 0)

	buttonTheme.Blurred.FocusedButton = buttonTheme.Blurred.FocusedButton.Background(lipgloss.Color("33")).Padding(1, 6).Border(lipgloss.NormalBorder(),false,false,false,false)
	buttonTheme.Focused.FocusedButton = buttonTheme.Focused.FocusedButton.Background(lipgloss.Color("33")).Padding(1, 6).Border(lipgloss.NormalBorder(), false, false, false, false)

	buttonTheme.Blurred.BlurredButton = buttonTheme.Blurred.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder(),false,false,false,false)
	buttonTheme.Focused.BlurredButton = buttonTheme.Focused.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder(), false, false, false, false)

	buttonTheme.Focused.Base = buttonTheme.Focused.Base.Border(lipgloss.NormalBorder(), false, false, false, false)
	buttonTheme.Blurred.Base = buttonTheme.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, false, false)

	var formData values

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("URL").
				Title("Enter Your Account URL").
				Value(&formData.url),

			huh.NewInput().
				Title("Enter Your Username").
				Key("USERNAME").
				Value(&formData.username),

			huh.NewInput().
				Title("Enter Your Password").
				Key("PASSWORD").
				Value(&formData.password).
				EchoMode(huh.EchoMode(textinput.EchoPassword)),

			huh.NewInput().
				Title("Enter Your Region").
				Key("REGION").
				Value(&formData.region),

			huh.NewInput().
				Title("Enter Your Tenant").
				Key("TENANT").
				Value(&formData.tenant),

			huh.NewConfirm().
				Key("CONFIRM").
				Title("Submit").
				Affirmative("Yes").
				Negative("No").
				Value(&formData.confirm).
				Validate(func(b bool) error {
					if !b {
						return nil
					}
					if formData.url == "" {
						return errors.New("enter valid url")
					}
					if formData.username == "" {
						return errors.New("enter valid username")
					}
					if formData.password == "" {
						return errors.New("enter valid password")
					}
					if formData.region == "" {
						return errors.New("enter valid region")
					}
					if formData.tenant == "" {
						return errors.New("enter valid tenant")
					}
					return nil
				}).WithTheme(buttonTheme),
		).Title("Configure pcdctl"),
	).WithTheme(t).WithShowHelp(false).WithWidth(70).WithKeyMap(&huh.KeyMap{
		Input: huh.InputKeyMap{
			Next: key.NewBinding(key.WithKeys("tab")),
			Prev: key.NewBinding(key.WithKeys("shift+tab")),
		},
		Confirm: huh.ConfirmKeyMap{
			Next:   key.NewBinding(key.WithKeys("tab")),
			Prev:   key.NewBinding(key.WithKeys("shift+tab")),
			Submit: key.NewBinding(key.WithKeys("enter")),
			Toggle: key.NewBinding(key.WithKeys("t")),
	},
	})

	return setup{
		state:    fillForm,
		form:     form,
		formData: formData,
	}
}
