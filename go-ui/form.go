package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)


func performInstllation(url, username, password, region, tenant string) error {

	if _,err:=exec.LookPath("pcdctl");err==nil{
		fmt.Println("pcdctl is alredy installed")
	} else {
		// Step 1: Download the script using curl and save locally
		scriptPath := "/tmp/pcdctl-setup.sh"
		curlCmd := exec.Command("curl", "-s", "-o", scriptPath, "https://pcdctl.s3.us-west-2.amazonaws.com/pcdctl-setup")
		curlCmd.Stdout = os.Stdout
		curlCmd.Stderr = os.Stderr

		if err := curlCmd.Run(); err != nil {
			return fmt.Errorf("failed to download pcdctl-setup script: %w", err)
		}

		// Step 2: Make the script executable
		if err := os.Chmod(scriptPath, 0755); err != nil {
			return fmt.Errorf("failed to chmod script: %w", err)
		}

		// Step 3: Run the script directly
		runCmd := exec.Command(scriptPath)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("failed to execute pcdctl setup script: %w", err)
		}
	}

	// Step 4: Check if pcdctl is already configured
	getCmd := exec.Command("pcdctl", "config", "get")
	getCmdOutput, err := getCmd.Output()

	if err != nil || len(getCmdOutput) == 0 {
		// Either error or empty config – proceed with setting config
		fmt.Println("ℹ️  Config not found. Running 'pcdctl config set'...")

		configCmd := exec.Command("pcdctl", "config", "set",
		"-u", url,
		"-e", username,
		"-p", password,
		"-r", region,
		"-t", tenant,
		)

		// following lines can be uncommented incase of Run instead of CombinedOuput
		// configCmd.Stdout = os.Stdout
		// configCmd.Stderr = os.Stderr

		if out,err := configCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to configure pcdctl: %w\nOutput:\n%s", err,string(out))
    		}
	} else {
		fmt.Println("✅ pcdctl already configured. Skipping config step.")
	}


	// Step 5: Prep the node
	prepCmd := exec.Command("pcdctl", "prep-node","-c","true")
	prepCmd.Stdout = os.Stdout
	prepCmd.Stderr = os.Stderr

	if err := prepCmd.Run(); err != nil {
		return fmt.Errorf("failed to run prep-node: %w", err)
	}

	fmt.Println("✅ PCD node setup complete!")
	return nil
}





type setup struct {
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

	form, cmd := s.form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		s.form = updated
	}

	if s.form.State == huh.StateCompleted {
		return s, func() tea.Msg {
			return configurationStatus("done")
		}
	}

	switch msg:=msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case "esc":
			return s, func() tea.Msg {
				return switchPageMsg("back")
			}
		}

	}
	
	return s, cmd
}

func (s setup) View() string {
	return s.form.View()
}

func FormInitialModel(width, height int) setup {
	t := huh.ThemeCharm()

	t.Focused.Base = t.Focused.Base.Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("39"))
	t.Blurred.Base = t.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, true, false)
	t.Blurred.Title = t.Blurred.Title.Foreground(lipgloss.Color("#3c3c3c")).Padding(0, 0, 1, 0)
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("39")).Padding(0, 0, 1, 0)
	t.Group.Title = t.Group.Title.Foreground(lipgloss.Color("#ffffff")).Bold(true).Padding(0, 0, 2, 0)
	t.Form.Base = t.Form.Base.Padding(3).AlignHorizontal(lipgloss.Center)

	buttonTheme := huh.ThemeCharm()

	buttonTheme.Blurred.Title = buttonTheme.Blurred.Title.Foreground(lipgloss.Color("39"))
	buttonTheme.Focused.Title = buttonTheme.Focused.Title.Foreground(lipgloss.Color("39"))

	buttonTheme.Blurred.Title = buttonTheme.Blurred.Title.Margin(0, 0, 1, 0)
	buttonTheme.Focused.Title = buttonTheme.Focused.Title.Margin(0, 0, 1, 0)

	buttonTheme.Blurred.FocusedButton = buttonTheme.Blurred.FocusedButton.Background(lipgloss.Color("0")).Padding(0, 7).Border(lipgloss.NormalBorder(),true,true,true,true)
	buttonTheme.Focused.FocusedButton = buttonTheme.Focused.FocusedButton.Background(lipgloss.Color("39")).Padding(1, 8).Border(lipgloss.NormalBorder(), false, false, false, false)

	buttonTheme.Blurred.BlurredButton = buttonTheme.Blurred.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder(),false,false,false,false)
	buttonTheme.Focused.BlurredButton = buttonTheme.Focused.BlurredButton.Padding(1, 6).Border(lipgloss.NormalBorder(), false, false, false, false)

	buttonTheme.Focused.Base = buttonTheme.Focused.Base.Border(lipgloss.NormalBorder(), false, false, false, false).AlignHorizontal(lipgloss.Center).Margin(1,0)
	buttonTheme.Blurred.Base = buttonTheme.Blurred.Base.Border(lipgloss.NormalBorder(), false, false, false, false).AlignHorizontal(lipgloss.Center).Margin(1,0)

	var formData values

	submitVal:=true

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
				Affirmative("Start pcdctl configuration").
				Negative("").
				Value(&submitVal).
				Validate(func(b bool) error {
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
					
					err:=performInstllation(formData.url,formData.username,formData.password,formData.region,formData.tenant)

					if err!=nil{
						return err
					}

					return nil

				}).WithTheme(buttonTheme),
		).Title("Configure pcdctl"),
	).WithTheme(t).WithShowHelp(false).WithHeight(int(0.7*float64(height))).WithWidth(int(0.5*float64(width))).WithKeyMap(&huh.KeyMap{
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
		form:     form,
		formData: formData,
	}
}
