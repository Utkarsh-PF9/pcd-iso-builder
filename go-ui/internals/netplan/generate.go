package netplan

import (
	"errors"
	"example/internals/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Network Network `yaml:"network"`
}

type Network struct {
	Version  int             `yaml:"version"`
	Renderer string          `yaml:"renderer"`
	Bonds    map[string]Bond `yaml:"bonds"`
}

type Bond struct {
	Interfaces  []string     `yaml:"interfaces"`
	Parameters  *Params      `yaml:"parameters,omitempty"`
	Addresses   []string     `yaml:"addresses,omitempty"`
	Gateway4    string       `yaml:"gateway4,omitempty"`
	Nameservers *Nameservers `yaml:"nameservers,omitempty"`
	DHCP4       bool         `yaml:"dhcp4,omitempty"`
}

type Params struct {
	Mode    string `yaml:"mode"`
	Primary string `yaml:"primary"`
}

type Nameservers struct {
	Addresses []string `yaml:"addresses"`
}




// to backup amy existing netplan.yaml files
// by renaming it to netplan.yaml.bak file
// backup existing netplan.yaml files (if /etc/netplan exists)
func backupNetplanConfigs() error {
	netplanDir := "/etc/netplan/"

	// Ensure the directory exists
	if err := os.MkdirAll(netplanDir, 0755); err != nil {
		return fmt.Errorf("failed to create netplan directory: %v", err)
	}

	// Proceed with renaming .yaml files
	return filepath.Walk(netplanDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			backupPath := path + ".bak"
			return os.Rename(path, backupPath)
		}
		return nil
	})
}




// generate and apply netplan according to the given values by the user
func GenerateAndApplyNetplan(data utils.Network_config_values) error {
	netplanDir := "/etc/netplan"

	// Ensure /etc/netplan exists
	if err := os.MkdirAll(netplanDir, 0755); err != nil {
		return fmt.Errorf("failed to create netplan directory: %v", err)
	}

	var network_config Config

	if !data.Is_dhcp {
		network_config = Config{
			Network: Network{
				Version:  2,
				Renderer: "networkd",
				Bonds: map[string]Bond{
					"bond0": {
						Interfaces: data.Interfaces,
						Parameters: &Params{
							Mode:    "active-backup",
							Primary: data.Interfaces[0],
						},
						Addresses: []string{fmt.Sprintf("%v/%v", data.Static_ip, data.CIDR)},
						Gateway4:  data.Gateway_ip,
						Nameservers: &Nameservers{
							Addresses: []string{data.Primary_dns, data.Secondary_dns},
						},
					},
				},
			},
		}
	} else {
		network_config = Config{
			Network: Network{
				Version:  2,
				Renderer: "networkd",
				Bonds: map[string]Bond{
					"bond0": {
						Interfaces: data.Interfaces,
						Parameters: &Params{
							Mode:    "active-backup",
							Primary: data.Interfaces[0],
						},
						DHCP4: true,
					},
				},
			},
		}
	}

	yamldata, err := yaml.Marshal(network_config)
	if err != nil {
		return errors.New("error generating netplan YAML content")
	}

	// Backup existing YAML files
	if err := backupNetplanConfigs(); err != nil {
		return err
	}

	// Write the new config
	configPath := filepath.Join(netplanDir, "pf9_netplan.yaml")
	if err := os.WriteFile(configPath, yamldata, 0644); err != nil {
		return errors.New("failed to write netplan config")
	}

	// Optionally apply it:
	// cmd := exec.Command("sudo", "netplan", "apply")
	// if output, err := cmd.CombinedOutput(); err != nil {
	//     return fmt.Errorf("netplan apply failed: %v\nOutput: %s", err, string(output))
	// }

	return nil
}

