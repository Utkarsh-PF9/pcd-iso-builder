#!/usr/bin/env python3

import curses
import subprocess
import json
import re
import sys
import time
from enum import Enum

class NetworkMode(Enum):
    DHCP = "dhcp"
    STATIC = "static"

class NetworkConfigUI:
    def __init__(self, stdscr):
        self.stdscr = stdscr
        self.current_field = 0
        self.network_mode = NetworkMode.DHCP
        
        # Network configuration data
        self.config = {
            'mode': NetworkMode.DHCP,
            'ip_address': '',
            'subnet_mask': '',
            'gateway': '',
            'dns_primary': '',
            'dns_secondary': '',
            'subnets': [
                {'name': 'mgmt', 'network': '192.168.10.0/24'},
                {'name': 'storage', 'network': '10.0.0.0/24'}
            ],
            'bond_device': '<bond0>',
            'bond_mode': 'active-backup',
            'interfaces': ['eth0'],
            'vlan_id': '',
            'ntp_server1': '',
            'ntp_server2': ''
        }
        
        # Field definitions for navigation
        self.fields = [
            'network_mode',
            'ip_address', 'subnet_mask', 'gateway',
            'dns_primary', 'dns_secondary',
            'subnets',
            'bonds',
            'vlan_id',
            'ntp_server1', 'ntp_server2'
        ]
        
        # Initialize curses
        curses.curs_set(0)  # Hide cursor
        self.stdscr.clear()
        self.height, self.width = self.stdscr.getmaxyx()
        
        # Color pairs
        curses.start_color()
        curses.init_pair(1, curses.COLOR_BLACK, curses.COLOR_CYAN)    # Header
        curses.init_pair(2, curses.COLOR_WHITE, curses.COLOR_BLUE)    # Selected field
        curses.init_pair(3, curses.COLOR_CYAN, curses.COLOR_BLACK)    # Labels
        curses.init_pair(4, curses.COLOR_YELLOW, curses.COLOR_BLACK)  # Values
        curses.init_pair(5, curses.COLOR_GREEN, curses.COLOR_BLACK)   # Active
        curses.init_pair(6, curses.COLOR_RED, curses.COLOR_BLACK)     # Inactive

    def draw_border(self):
        """Draw the main border and title"""
        self.stdscr.border()
        title = "Private Cloud Director - Network Setup"
        title_x = (self.width - len(title)) // 2
        self.stdscr.addstr(0, title_x, title, curses.color_pair(1) | curses.A_BOLD)

    def draw_network_mode(self, y_start):
        """Draw network mode selection"""
        y = y_start
        self.stdscr.addstr(y, 4, "Network mode:", curses.color_pair(3))
        
        # DHCP option
        dhcp_attrs = curses.color_pair(5) if self.config['mode'] == NetworkMode.DHCP else curses.color_pair(6)
        dhcp_marker = "●" if self.config['mode'] == NetworkMode.DHCP else "○"
        self.stdscr.addstr(y, 20, f"( {dhcp_marker} ) DHCP", dhcp_attrs)
        
        # Static IP option  
        static_attrs = curses.color_pair(5) if self.config['mode'] == NetworkMode.STATIC else curses.color_pair(6)
        static_marker = "●" if self.config['mode'] == NetworkMode.STATIC else "○"
        self.stdscr.addstr(y, 45, f"( {static_marker} ) Static IP", static_attrs)
        
        return y + 2

    def draw_static_details(self, y_start):
        """Draw static IP configuration section"""
        y = y_start
        self.stdscr.addstr(y, 4, "── Static-IP details ──────────────────────", curses.color_pair(3))
        y += 1
        
        # IP Address
        self.stdscr.addstr(y, 6, "IP address  :", curses.color_pair(3))
        ip_attrs = curses.color_pair(2) if self.current_field == 1 else curses.color_pair(4)
        ip_value = self.config['ip_address'].ljust(15)
        self.stdscr.addstr(y, 20, f"[{ip_value}]", ip_attrs)
        y += 1
        
        # Subnet mask
        self.stdscr.addstr(y, 6, "Subnet mask :", curses.color_pair(3))
        mask_attrs = curses.color_pair(2) if self.current_field == 2 else curses.color_pair(4)
        mask_value = self.config['subnet_mask'].ljust(15)
        self.stdscr.addstr(y, 20, f"[{mask_value}]", mask_attrs)
        y += 1
        
        # Gateway
        self.stdscr.addstr(y, 6, "Gateway     :", curses.color_pair(3))
        gw_attrs = curses.color_pair(2) if self.current_field == 3 else curses.color_pair(4)
        gw_value = self.config['gateway'].ljust(15)
        self.stdscr.addstr(y, 20, f"[{gw_value}]", gw_attrs)
        
        return y + 2

    def draw_subnets(self, y_start):
        """Draw subnets section"""
        y = y_start
        self.stdscr.addstr(y, 4, "Subnet(s) (1/1 to reorder, SPACE=toggle use)", curses.color_pair(3))
        y += 1
        
        for i, subnet in enumerate(self.config['subnets']):
            marker = "x" if i == 0 else " "
            self.stdscr.addstr(y + i, 6, f"[{marker}] {subnet['network']:<20} {subnet['name']:<10}", curses.color_pair(4))
            if subnet['name'] == 'mgmt':
                self.stdscr.addstr(y + i, 50, "[Add] [Remove]", curses.color_pair(3))
            else:
                self.stdscr.addstr(y + i, 50, "[Add] [Remove]", curses.color_pair(3))
        
        return y + len(self.config['subnets']) + 1

    def draw_dns_servers(self, y_start):
        """Draw DNS servers section"""
        y = y_start
        self.stdscr.addstr(y, 4, "DNS servers", curses.color_pair(3))
        y += 1
        
        # Primary DNS
        self.stdscr.addstr(y, 6, "Primary     :", curses.color_pair(3))
        pri_attrs = curses.color_pair(2) if self.current_field == 4 else curses.color_pair(4)
        pri_value = self.config['dns_primary'].ljust(15)
        self.stdscr.addstr(y, 20, f"[{pri_value}]", pri_attrs)
        y += 1
        
        # Secondary DNS
        self.stdscr.addstr(y, 6, "Secondary   :", curses.color_pair(3))
        sec_attrs = curses.color_pair(2) if self.current_field == 5 else curses.color_pair(4)
        sec_value = self.config['dns_secondary'].ljust(15)
        self.stdscr.addstr(y, 20, f"[{sec_value}]", sec_attrs)
        
        return y + 2

    def draw_bonds_teams(self, y_start):
        """Draw bonds/teams section"""
        y = y_start
        self.stdscr.addstr(y, 45, "Bonds / Teams", curses.color_pair(3))
        y += 1
        self.stdscr.addstr(y, 45, f"Bond device : {self.config['bond_device']} ▼", curses.color_pair(4))
        y += 1
        self.stdscr.addstr(y, 45, f"Mode        : {self.config['bond_mode']} ▼", curses.color_pair(4))
        y += 1
        self.stdscr.addstr(y, 45, f"Interfaces  : [eth0] [eth1]", curses.color_pair(4))
        
        return y + 2

    def draw_vlan(self, y_start):
        """Draw VLAN configuration"""
        y = y_start
        self.stdscr.addstr(y, 4, "VLAN tagging", curses.color_pair(3))
        y += 1
        self.stdscr.addstr(y, 6, "[x] Enable VLANs", curses.color_pair(4))
        
        vlan_attrs = curses.color_pair(2) if self.current_field == 8 else curses.color_pair(4)
        vlan_value = self.config['vlan_id'].ljust(10)
        self.stdscr.addstr(y, 45, f"VLAN ID : [{vlan_value}] Interface : <eth0 ▼>", vlan_attrs)
        
        return y + 2

    def draw_ntp(self, y_start):
        """Draw NTP configuration"""
        y = y_start
        self.stdscr.addstr(y, 4, "NTP configuration", curses.color_pair(3))
        y += 1
        
        # NTP Server 1
        self.stdscr.addstr(y, 6, "Server 1    :", curses.color_pair(3))
        ntp1_attrs = curses.color_pair(2) if self.current_field == 9 else curses.color_pair(4)
        ntp1_value = self.config['ntp_server1'].ljust(25)
        self.stdscr.addstr(y, 20, f"[{ntp1_value}]", ntp1_attrs)
        y += 1
        
        # NTP Server 2
        self.stdscr.addstr(y, 6, "Server 2    :", curses.color_pair(3))
        ntp2_attrs = curses.color_pair(2) if self.current_field == 10 else curses.color_pair(4)
        ntp2_value = self.config['ntp_server2'].ljust(25)
        self.stdscr.addstr(y, 20, f"[{ntp2_value}]", ntp2_attrs)
        
        return y + 2

    def draw_help(self):
        """Draw help information at bottom"""
        help_text = "TAB: Next field | SPACE: Toggle | ENTER: Edit | ESC: Exit"
        self.stdscr.addstr(self.height - 2, 2, help_text, curses.color_pair(1))

    def draw_screen(self):
        """Main screen drawing function"""
        self.stdscr.clear()
        self.draw_border()
        
        y = 2
        y = self.draw_network_mode(y)
        y = self.draw_static_details(y)
        y = self.draw_subnets(y)
        y = self.draw_dns_servers(y)
        
        # Right column
        y_right = 8
        y_right = self.draw_bonds_teams(y_right)
        
        y = max(y, y_right) + 1
        y = self.draw_vlan(y)
        y = self.draw_ntp(y)
        
        self.draw_help()
        self.stdscr.refresh()

    def handle_input(self, key):
        """Handle keyboard input"""
        if key == 27:  # ESC
            return False
        elif key == ord('\t') or key == curses.KEY_DOWN:  # TAB or Down
            self.current_field = (self.current_field + 1) % len(self.fields)
        elif key == curses.KEY_UP:
            self.current_field = (self.current_field - 1) % len(self.fields)
        elif key == ord(' '):  # SPACE - toggle
            if self.current_field == 0:  # Network mode
                self.config['mode'] = NetworkMode.STATIC if self.config['mode'] == NetworkMode.DHCP else NetworkMode.DHCP
        elif key == ord('\n') or key == curses.KEY_ENTER or key == 10:  # ENTER - edit
            self.edit_field()
        
        return True

    def edit_field(self):
        """Edit the current field"""
        field_map = {
            1: 'ip_address',
            2: 'subnet_mask', 
            3: 'gateway',
            4: 'dns_primary',
            5: 'dns_secondary',
            8: 'vlan_id',
            9: 'ntp_server1',
            10: 'ntp_server2'
        }
        
        if self.current_field in field_map:
            field_name = field_map[self.current_field]
            current_value = self.config[field_name]
            
            # Create edit window
            edit_win = curses.newwin(3, 40, self.height//2 - 1, self.width//2 - 20)
            edit_win.border()
            edit_win.addstr(1, 2, f"Edit {field_name}: ")
            edit_win.addstr(1, 15, current_value)
            edit_win.refresh()
            
            # Simple edit (in real implementation you'd want proper text input)
            curses.echo()
            curses.curs_set(1)
            new_value = edit_win.getstr(1, 15, 20).decode('utf-8')
            curses.noecho()
            curses.curs_set(0)
            
            self.config[field_name] = new_value
            del edit_win

    def run(self):
        """Main run loop"""
        while True:
            self.draw_screen()
            key = self.stdscr.getch()
            if not self.handle_input(key):
                break

def show_system_menu(stdscr):
    """Show system management menu"""
    menu_items = [
        "Network Configuration",
        "System Information", 
        "Service Management",
        "Disk Management",
        "Log Viewer",
        "Exit to Shell"
    ]
    
    current_item = 0
    
    while True:
        stdscr.clear()
        h, w = stdscr.getmaxyx()
        
        # Title
        title = "Ubuntu System Management"
        stdscr.addstr(1, (w - len(title)) // 2, title, curses.color_pair(1) | curses.A_BOLD)
        
        # Menu items
        for i, item in enumerate(menu_items):
            if i == current_item:
                stdscr.addstr(5 + i, (w - len(item)) // 2, item, curses.color_pair(2) | curses.A_BOLD)
            else:
                stdscr.addstr(5 + i, (w - len(item)) // 2, item, curses.color_pair(3))
        
        stdscr.addstr(h - 2, 2, "↑↓: Navigate | ENTER: Select | ESC: Exit", curses.color_pair(1))
        stdscr.refresh()
        
        key = stdscr.getch()
        if key == curses.KEY_UP and current_item > 0:
            current_item -= 1
        elif key == curses.KEY_DOWN and current_item < len(menu_items) - 1:
            current_item += 1
        elif key == ord('\n') or key == 10:
            if current_item == 0:  # Network Configuration
                network_ui = NetworkConfigUI(stdscr)
                network_ui.run()
            elif current_item == len(menu_items) - 1:  # Exit to Shell
                break
            else:
                # Placeholder for other menu items
                stdscr.addstr(h//2, (w - 20)//2, "Feature coming soon!", curses.color_pair(4))
                stdscr.addstr(h//2 + 1, (w - 25)//2, "Press any key to continue", curses.color_pair(3))
                stdscr.refresh()
                stdscr.getch()
        elif key == 27:  # ESC
            break

def main():
    try:
        curses.wrapper(show_system_menu)
        print("Exiting to shell...")
        print("Type 'custom-ui' to return to the management interface.")
    except KeyboardInterrupt:
        print("\nExiting...")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main() 