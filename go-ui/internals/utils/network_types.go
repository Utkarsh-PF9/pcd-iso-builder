package utils

type Network_config_values struct {
	Is_dhcp       bool
	Static_ip     string
	CIDR          string
	Gateway_ip    string
	Primary_dns   string
	Secondary_dns string
	Interfaces    []string
}
