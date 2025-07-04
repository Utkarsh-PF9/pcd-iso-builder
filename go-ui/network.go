package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	netplan "example/internals/netplan"
	utils "example/internals/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func isValidIPv4(ip string) bool {
	return net.ParseIP(ip) != nil
}

type network struct {
	network_form *huh.Form
	form_data    utils.Network_config_values
}

var bonds_list_items = []list.Item{
	item{title: "Bond 0", desc: "Not configured"},
}

func (n network) Init() tea.Cmd {
	return n.network_form.Init()
}
func (n network) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return n, func() tea.Msg {
				return switchPageMsg("back")
			}
		}
	}

	form, cmd := n.network_form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		n.network_form = updated

		err:=n.network_form.Errors()

		if n.network_form.State==huh.StateCompleted && len(err)==0{
			return n, func() tea.Msg {
				return switchPageMsg("back")
			}
		}
	}

	return n, cmd
}

func (n network) View() string {
	return n.network_form.View()
}

func NetworkInitialModel() network {

	t := huh.ThemeCharm()

	t.Focused.Base = t.Focused.Base.Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("33"))
	t.Blurred.Base = t.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, true, false)

	t.Blurred.Title = t.Blurred.Title.Foreground(lipgloss.Color("#3c3c3c")).Padding(0, 0, 1, 0)

	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("33")).Padding(0, 0, 1, 0)

	t.Blurred.Description = t.Blurred.Description.Margin(1, 0)
	t.Focused.Description = t.Focused.Description.Margin(0, 0, 1, 0)

	t.Group.Title = t.Group.Title.Foreground(lipgloss.Color("#ffffff")).Bold(true).Padding(0, 0, 1, 0)
	t.Form.Base = t.Form.Base.Padding(1).AlignHorizontal(lipgloss.Center)

	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(lipgloss.Color("33")).Margin(2, 0, 0, 0)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color("33")).Margin(2, 0, 0, 0)

	selectTheme := huh.ThemeCharm()

	selectTheme.Focused.Base = selectTheme.Focused.Base.Border(lipgloss.NormalBorder(), false, false, false, false).Margin(0, 0, 1, 0)
	selectTheme.Blurred.Base = selectTheme.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, false, false).Margin(0, 0, 1, 0)

	selectTheme.Blurred.Title = selectTheme.Blurred.Title.Foreground(lipgloss.Color("33")).Margin(0, 0, 1, 0)
	selectTheme.Focused.Title = selectTheme.Focused.Title.Foreground(lipgloss.Color("33")).Margin(0, 0, 1, 0)

	selectTheme.Blurred.UnselectedOption = selectTheme.Blurred.UnselectedOption.Foreground(lipgloss.Color("0"))
	selectTheme.Focused.UnselectedOption = selectTheme.Focused.UnselectedOption.Foreground(lipgloss.Color("0"))

	selectTheme.Blurred.SelectedOption = selectTheme.Blurred.SelectedOption.Foreground(lipgloss.Color("#ffffff"))
	selectTheme.Focused.SelectedOption = selectTheme.Focused.SelectedOption.Foreground(lipgloss.Color("#ffffff"))

	// unused
	// t.Form.Base=t.Form.Base.Width(100)
	// t.Form.Base = t.Form.Base.Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("#ffffff")).Padding(3).AlignHorizontal(lipgloss.Center)

	// t.Blurred.TextInput.Placeholder = t.Blurred.TextInput.Placeholder.Background(lipgloss.Color("#ffffff"))
	// t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("#03add7")).Bold(true)
	// unused end

	button_t := huh.ThemeCharm()

	button_t.Blurred.Title = button_t.Blurred.Title.Foreground(lipgloss.Color("33"))
	button_t.Focused.Title = button_t.Focused.Title.Foreground(lipgloss.Color("33"))

	button_t.Blurred.Title = button_t.Blurred.Title.Margin(0, 0, 1, 0)
	button_t.Focused.Title = button_t.Focused.Title.Margin(0, 0, 1, 0)

	button_t.Blurred.FocusedButton = button_t.Blurred.FocusedButton.Background(lipgloss.Color("0")).Padding(1, 6).Border(lipgloss.NormalBorder())
	button_t.Focused.FocusedButton = button_t.Focused.FocusedButton.Background(lipgloss.Color("33")).Padding(1, 6).Border(lipgloss.NormalBorder(), false, false, false, false).Margin(0, 1)

	button_t.Blurred.BlurredButton = button_t.Blurred.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder())
	button_t.Focused.BlurredButton = button_t.Focused.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder(), false, false, false, false).Margin(0, 1)

	button_t.Focused.Base = button_t.Focused.Base.Border(lipgloss.NormalBorder(), false, false, false, false)
	button_t.Blurred.Base = button_t.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, false, false)

	// button_t.Blurred.SelectedOption=button_t.Blurred.SelectedOption.Border(lipgloss.NormalBorder(),true,true,true,true).Background(lipgloss.Color("0")).Foreground(lipgloss.Color("#ffffff")).Padding(1,5).Width(50).Align(lipgloss.Center)
	// button_t.Blurred.UnselectedOption=button_t.Blurred.UnselectedOption.Border(lipgloss.NormalBorder(),true,true,true,true).Background(lipgloss.Color("0")).Foreground(lipgloss.Color("#ffffff")).Padding(1,5).Width(50).Align(lipgloss.Center)

	// button_t.Focused.SelectedOption=button_t.Focused.SelectedOption.Background(lipgloss.Color("33")).Foreground(lipgloss.Color("#ffffff")).Padding(1,5).Margin(1).Width(50).Align(lipgloss.Center)
	// button_t.Focused.UnselectedOption=button_t.Focused.UnselectedOption.Background(lipgloss.Color("33")).Foreground(lipgloss.Color("#ffffff")).Padding(1,5).Margin(1).Width(50).Align(lipgloss.Center)

	// button_t.Focused.Base = button_t.Focused.Base.Border(lipgloss.NormalBorder(), true, true, true, true).BorderForeground(lipgloss.Color("33"))
	// button_t.Blurred.Base = button_t.Blurred.Base.Border(lipgloss.NormalBorder(), true, true, true, true)

	// button_t.Blurred.Title = button_t.Blurred.Title.Foreground(lipgloss.Color("33"))

	// button_t.Focused.Title = button_t.Focused.Title.Foreground(lipgloss.Color("33"))

	// button_t.Blurred.BlurredButton = button_t.Blurred.BlurredButton.Padding(1, 6)
	// button_t.Focused.BlurredButton = button_t.Focused.BlurredButton.Padding(1, 6)

	// button_t.Blurred.FocusedButton = button_t.Blurred.FocusedButton.Background(lipgloss.Color("33")).Padding(1, 6)
	// button_t.Focused.FocusedButton = button_t.Focused.FocusedButton.Background(lipgloss.Color("33")).Padding(1, 6)

	// button_t.Focused.Base = button_t.Focused.Base.Border(lipgloss.NormalBorder(), false, false, false, false)
	// button_t.Blurred.Base = button_t.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, false, false)

	bond_theme := huh.ThemeCharm()
	// bond_theme.Blurred.SelectedOption=bond_theme.Blurred.SelectedOption.Foreground(lipgloss.Color("33")).Border(lipgloss.NormalBorder())
	// bond_theme.Blurred.Option=bond_theme.Blurred.Option.Foreground(lipgloss.Color("33"))
	// bond_theme.Focused.Option=bond_theme.Focused.Option.Foreground(lipgloss.Color("33"))

	// bond_theme.Blurred.SelectSelector=bond_theme.Blurred.SelectSelector.Foreground(lipgloss.Color("33"))

	// bond_theme.Blurred.TextInput.Placeholder=bond_theme.Blurred.TextInput.Placeholder.Foreground(lipgloss.Color("33"))

	// bond_theme.Blurred.TextInput.Text=bond_theme.Blurred.TextInput.Text.Foreground(lipgloss.Color("33"))

	bond_theme.Focused.Base = bond_theme.Focused.Base.Border(lipgloss.NormalBorder(), false, false, true, false)
	bond_theme.Blurred.Base = bond_theme.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, true, false)

	bond_theme.Blurred.Title = bond_theme.Blurred.Title.Foreground(lipgloss.Color("#3c3c3c")).Padding(0, 0, 1, 0)
	bond_theme.Focused.Title = bond_theme.Focused.Title.Foreground(lipgloss.Color("33")).Padding(0, 0, 1, 0)
	bond_theme.Blurred.Description = bond_theme.Blurred.Description.Margin(1, 0)
	bond_theme.Focused.Description = bond_theme.Focused.Description.Margin(0, 0, 1, 0)

	bond_theme.Blurred.Option = bond_theme.Blurred.Option.Margin(0, 0, 1, 0)
	bond_theme.Blurred.UnselectedOption = bond_theme.Blurred.UnselectedOption.Margin(0, 0, 1, 0)
	bond_theme.Blurred.SelectedOption = bond_theme.Blurred.SelectedOption.Margin(0, 0, 1, 0)

	bond_theme.Focused.Option = bond_theme.Focused.Option.Margin(0, 0, 1, 0)
	bond_theme.Focused.UnselectedOption = bond_theme.Focused.UnselectedOption.Margin(0, 0, 1, 0)
	bond_theme.Focused.SelectedPrefix = bond_theme.Focused.SelectedPrefix.Foreground(lipgloss.Color("33"))
	bond_theme.Focused.SelectedOption = bond_theme.Focused.SelectedOption.Margin(0, 0, 1, 0).Foreground(lipgloss.Color("33"))

	bond_theme.Blurred.Description = bond_theme.Blurred.Description.Foreground(lipgloss.Color("0")).Margin(0, 0, 1, 0)
	bond_theme.Focused.Description = bond_theme.Focused.Description.Foreground(lipgloss.Color("#ffffff")).Margin(0, 0, 1, 0)
	bond_theme.Blurred.SelectedPrefix = bond_theme.Blurred.SelectedPrefix.Foreground(lipgloss.Color("33"))
	bond_theme.Blurred.SelectedOption = bond_theme.Blurred.SelectedOption.Foreground(lipgloss.Color("33"))
	bond_theme.Blurred.UnselectedOption = bond_theme.Blurred.UnselectedOption.Foreground(lipgloss.Color("0"))

	// bond_theme.Focused.SelectedPrefix=bond_theme.Focused.SelectedPrefix.SetString("")
	// bond_theme.Blurred.SelectedPrefix=bond_theme.Blurred.SelectedPrefix.SetString("")

	bond_theme.Focused.UnselectedPrefix = bond_theme.Focused.UnselectedPrefix.SetString("  ")
	bond_theme.Blurred.UnselectedPrefix = bond_theme.Blurred.UnselectedPrefix.SetString("  ")

	// someBool := false

	var options []huh.Option[string]

	ifaces, _ := net.Interfaces()
	for i, iface := range ifaces {
		addrs, _ := iface.Addrs()
		var ipList []string
		for _, addr := range addrs {
			ipList = append(ipList, addr.String())
		}
		ipStr := "None"
		if len(ipList) > 0 {
			ipStr = fmt.Sprintf("%v", ipList)
		}

		info := fmt.Sprintf(
			"%d. %s\n   MTU: %d | MAC: %s\n   Flags: %v\n   IPs: %s\n",
			i+1, iface.Name, iface.MTU, iface.HardwareAddr.String(), iface.Flags.String(), ipStr,
		)

		// info := fmt.Sprintf("%d. %v", i+1, iface.Name)

		options = append(options, huh.NewOption(info, iface.Name))
	}

	saveBool := true

	form_data_values := utils.Network_config_values{
		Is_dhcp:       true,
		Static_ip:     "",
		CIDR:   "",
		Gateway_ip:    "",
		Primary_dns:   "8.8.8.8",
		Secondary_dns: "1.1.1.1",
		Interfaces:    []string{},
	}

	return network{
		network_form: huh.NewForm(
			huh.NewGroup(
				// huh.NewSelect[string]().
				// 	Title("Select IP Type").
				// 	Options(
				// 		huh.Option[string]{
				// 			Key:   "Static",
				// 			Value: "Static",
				// 		},
				// 		huh.Option[string]{
				// 			Key:   "DHCP",
				// 			Value: "DHCP",
				// 		},
				// 	).Value(&selectedType).WithTheme(selectTheme).WithKeyMap(&huh.KeyMap{
				// 		Select: huh.SelectKeyMap{
				// 			Up: key.NewBinding(key.WithKeys("t")),
				// 			Down: key.NewBinding(key.WithKeys("t")),
				// 		},
				// 	}),

				// huh.NewNote().Title("Select a Method:"),
				huh.NewConfirm().
					Title("Select a Method:").
					Description("toggle: t").
					Affirmative("DHCP").
					Negative("Static IP").
					Value(&form_data_values.Is_dhcp).
					WithTheme(button_t),

				huh.NewNote().Title("Static IP Configuration:").Description("Note: optional in case of DHCP"),
				huh.NewInput().Title("Enter Static IP").Placeholder("Ex: 192.168.22.23").
					Value(&form_data_values.Static_ip),

				huh.NewInput().Title("Enter CIDR").Placeholder("Ex: 16").
					Value(&form_data_values.CIDR),

				huh.NewInput().Title("Enter Gateway IP").Placeholder("Ex: 192.168.0.1").
					Value(&form_data_values.Gateway_ip),


				huh.NewNote().Title("DNS:").Description("Note: optional in case of DHCP"),
				huh.NewInput().Title("Primary:").Placeholder("Ex: 8.8.8.8").Value(&form_data_values.Primary_dns),
				huh.NewInput().Title("Secondary:").Placeholder("Ex: 1.1.1.1").Value(&form_data_values.Secondary_dns),

				huh.NewMultiSelect[string]().
					Key("INTERFACES").
					Title("Select at least one interface for Bond 0:").
					Description("move: ↑/↓  •  select: x  •  select all: shift + a  •  select none: shift + n").
					Options(options...).
					Value(&form_data_values.Interfaces).WithTheme(bond_theme),

				huh.NewConfirm().
					Title("Save").
					Description("toggle: t").
					Affirmative("Yes").
					Negative("No").
					Value(&saveBool).
					Validate(func(b bool) error {
						if !b {
							return nil
						}

						if !b {
							return errors.New("you must confirm to save and proceed")
						}

						if form_data_values.Is_dhcp {
							if len(form_data_values.Interfaces) == 0 {
								return errors.New("please select at least one network interface")
							}

							if !isValidIPv4(form_data_values.Primary_dns) {
								form_data_values.Primary_dns = "8.8.8.8"
							}

							if !isValidIPv4(form_data_values.Secondary_dns) {
								form_data_values.Secondary_dns = "1.1.1.1"
							}

							return netplan.GenerateAndApplyNetplan(form_data_values)
						}

						i_CIDR,_:=strconv.Atoi(form_data_values.CIDR)
						if i_CIDR>32 || i_CIDR==0 {
							return errors.New("invalid CIDR, must be between 1 and 32")
						}

						// Static IP path
						if !isValidIPv4(form_data_values.Static_ip) {
							return errors.New("invalid static IP address")
						}


						if !isValidIPv4(form_data_values.Gateway_ip) {
							return errors.New("invalid gateway IP")
						}

						if !isValidIPv4(form_data_values.Primary_dns) {
							return errors.New("invalid primary DNS IP")
						}

						if !isValidIPv4(form_data_values.Secondary_dns) {
							return errors.New("invalid secondary DNS IP")
						}

						return netplan.GenerateAndApplyNetplan(form_data_values)
					}).
					WithTheme(button_t),
			),
		).WithKeyMap(&huh.KeyMap{
			Input: huh.InputKeyMap{
				Next: key.NewBinding(key.WithKeys("tab", "enter")),
				Prev: key.NewBinding(key.WithKeys("shift+tab")),
			},
			Confirm: huh.ConfirmKeyMap{
				Toggle: key.NewBinding(key.WithKeys("t")),
				Prev:   key.NewBinding(key.WithKeys("shift+tab")),
				Next:   key.NewBinding(key.WithKeys("tab")),
				Submit: key.NewBinding(key.WithKeys("enter")),
			},
			MultiSelect: huh.MultiSelectKeyMap{
				Next:       key.NewBinding(key.WithKeys("tab")),
				Prev:       key.NewBinding(key.WithKeys("shift+tab")),
				Up:         key.NewBinding(key.WithKeys("up")),
				Down:       key.NewBinding(key.WithKeys("down")),
				Toggle:     key.NewBinding(key.WithKeys("x")),
				Submit:     key.NewBinding(key.WithKeys("enter")),
				SelectAll:  key.NewBinding(key.WithKeys("A")),
				SelectNone: key.NewBinding(key.WithKeys("N")),
			},
		}).WithHeight(40).WithWidth(80).WithTheme(t),

		form_data: form_data_values,
	}
}
