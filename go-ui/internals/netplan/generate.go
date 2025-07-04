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
func backupNetplanConfigs() error {
	netplanDir := "/etc/netplan/"
	return filepath.Walk(netplanDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Process only .yaml files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			backupPath := path + ".bak"
			return os.Rename(path, backupPath)
		}
		return nil
	})
}

// generate and apply netplan according to the given values by the user
func GenerateAndApplyNetplan(data utils.Network_config_values) error {
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
							Primary: "eth0",
						},
						DHCP4: true,
					},
				},
			},
		}
	}

	yamldata, err := yaml.Marshal(network_config)

	if err != nil {
		return errors.New("error generating netplan")
	}

	err = backupNetplanConfigs()
	if err != nil {
		return errors.New("error while backing up netplan yaml files in /etc/netplan")
	}

	err = os.WriteFile("/etc/netplan/pf9_netplan.yaml", yamldata, 0644)
	if err != nil {
		return errors.New("error creating netplan.yaml in /etc/netplan")
	}

	// cmd := exec.Command("sudo", "netplan", "apply")
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return fmt.Errorf("netplan apply failed: %v\nOutput: %s", err, string(output))
	// }
	return nil

}
