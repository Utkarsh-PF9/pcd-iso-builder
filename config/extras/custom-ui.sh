#!/bin/bash

# Colors for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Function to clear screen and show header
show_header() {
    clear
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${WHITE}                    Ubuntu System Management                   ${GREEN}║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo
}

# Function to show main menu
show_menu() {
    show_header
    echo -e "${YELLOW}Please select an option:${NC}"
    echo
    echo -e "${CYAN}[1]${NC} System Information"
    echo -e "${CYAN}[2]${NC} Network Status"
    echo -e "${CYAN}[3]${NC} Disk Usage"
    echo -e "${CYAN}[4]${NC} Running Services"
    echo -e "${CYAN}[5]${NC} Memory Usage"
    echo -e "${CYAN}[6]${NC} Process Monitor"
    echo -e "${CYAN}[7]${NC} System Logs"
    echo -e "${CYAN}[0]${NC} Exit"
    echo
    echo -e "${WHITE}Enter your choice [0-7]: ${NC}"
}

# Function to wait for user input
wait_for_key() {
    echo
    echo -e "${YELLOW}Press any key to continue...${NC}"
    read -n 1 -s
}

# Function to show system information
show_system_info() {
    show_header
    echo -e "${BLUE}System Information:${NC}"
    echo -e "${WHITE}==================${NC}"
    echo
    
    echo -e "${CYAN}Hostname:${NC} $(hostname)"
    echo -e "${CYAN}Kernel:${NC} $(uname -r)"
    echo -e "${CYAN}Architecture:${NC} $(uname -m)"
    echo -e "${CYAN}OS Release:${NC}"
    if [ -f /etc/os-release ]; then
        source /etc/os-release
        echo "  $PRETTY_NAME"
    fi
    echo
    
    echo -e "${CYAN}Uptime:${NC}"
    uptime
    echo
    
    echo -e "${CYAN}Load Average:${NC}"
    cat /proc/loadavg
    echo
    
    echo -e "${CYAN}CPU Info:${NC}"
    echo "  $(grep 'model name' /proc/cpuinfo | head -1 | cut -d: -f2 | sed 's/^ *//')"
    echo "  Cores: $(nproc)"
    
    wait_for_key
}

# Function to show network status
show_network_status() {
    show_header
    echo -e "${BLUE}Network Status:${NC}"
    echo -e "${WHITE}===============${NC}"
    echo
    
    echo -e "${CYAN}Active Network Interfaces:${NC}"
    ip -brief addr show | grep -v "lo\|UNKNOWN" | while read line; do
        echo "  $line"
    done
    echo
    
    echo -e "${CYAN}Routing Table:${NC}"
    ip route | head -5
    echo
    
    echo -e "${CYAN}DNS Servers:${NC}"
    if [ -f /etc/resolv.conf ]; then
        grep nameserver /etc/resolv.conf | head -3
    fi
    echo
    
    echo -e "${CYAN}Network Connections:${NC}"
    ss -tuln | head -10
    
    wait_for_key
}

# Function to show disk usage
show_disk_usage() {
    show_header
    echo -e "${BLUE}Disk Usage:${NC}"
    echo -e "${WHITE}============${NC}"
    echo
    
    echo -e "${CYAN}Filesystem Usage:${NC}"
    df -h | grep -vE '^Filesystem|tmpfs|cdrom|udev'
    echo
    
    echo -e "${CYAN}Disk I/O Statistics:${NC}"
    if command -v iostat >/dev/null 2>&1; then
        iostat -x 1 1 | tail -n +4
    else
        echo "  iostat not available (install sysstat package)"
    fi
    echo
    
    echo -e "${CYAN}Largest Directories in /home:${NC}"
    if [ -d /home ]; then
        du -sh /home/* 2>/dev/null | sort -hr | head -5
    fi
    
    wait_for_key
}

# Function to show running services
show_services() {
    show_header
    echo -e "${BLUE}Running Services:${NC}"
    echo -e "${WHITE}=================${NC}"
    echo
    
    echo -e "${CYAN}Active Services (systemd):${NC}"
    systemctl list-units --type=service --state=running --no-pager --no-legend | head -15 | while read line; do
        service_name=$(echo "$line" | awk '{print $1}')
        status=$(echo "$line" | awk '{print $3}')
        echo -e "  ${GREEN}●${NC} $service_name - $status"
    done
    echo
    
    echo -e "${CYAN}Failed Services:${NC}"
    failed_services=$(systemctl list-units --type=service --state=failed --no-pager --no-legend)
    if [ -z "$failed_services" ]; then
        echo -e "  ${GREEN}No failed services${NC}"
    else
        echo "$failed_services" | head -5 | while read line; do
            service_name=$(echo "$line" | awk '{print $1}')
            echo -e "  ${RED}●${NC} $service_name"
        done
    fi
    
    wait_for_key
}

# Function to show memory usage
show_memory_usage() {
    show_header
    echo -e "${BLUE}Memory Usage:${NC}"
    echo -e "${WHITE}==============${NC}"
    echo
    
    echo -e "${CYAN}Memory Information:${NC}"
    free -h
    echo
    
    echo -e "${CYAN}Memory Usage by Process (Top 10):${NC}"
    ps aux --sort=-%mem | head -11 | tail -10 | while read line; do
        user=$(echo "$line" | awk '{print $1}')
        pid=$(echo "$line" | awk '{print $2}')
        mem=$(echo "$line" | awk '{print $4}')
        cmd=$(echo "$line" | awk '{for(i=11;i<=NF;i++) printf "%s ", $i; print ""}' | cut -c1-40)
        printf "  %-10s %6s %5s%% %s\n" "$user" "$pid" "$mem" "$cmd"
    done
    echo
    
    echo -e "${CYAN}Swap Usage:${NC}"
    swapon --show 2>/dev/null || echo "  No swap configured"
    
    wait_for_key
}

# Function to show process monitor
show_processes() {
    show_header
    echo -e "${BLUE}Process Monitor:${NC}"
    echo -e "${WHITE}=================${NC}"
    echo
    
    echo -e "${CYAN}Top Processes by CPU Usage:${NC}"
    ps aux --sort=-%cpu | head -11 | tail -10 | while read line; do
        user=$(echo "$line" | awk '{print $1}')
        pid=$(echo "$line" | awk '{print $2}')
        cpu=$(echo "$line" | awk '{print $3}')
        cmd=$(echo "$line" | awk '{for(i=11;i<=NF;i++) printf "%s ", $i; print ""}' | cut -c1-40)
        printf "  %-10s %6s %5s%% %s\n" "$user" "$pid" "$cpu" "$cmd"
    done
    echo
    
    echo -e "${CYAN}Process Count by User:${NC}"
    ps aux | awk 'NR>1 {count[$1]++} END {for (user in count) printf "  %-15s %d\n", user, count[user]}' | sort -k2 -nr | head -5
    echo
    
    echo -e "${CYAN}Zombie Processes:${NC}"
    zombie_count=$(ps aux | awk '$8 ~ /^Z/ { count++ } END { print count+0 }')
    echo "  Count: $zombie_count"
    
    wait_for_key
}

# Function to show system logs
show_logs() {
    show_header
    echo -e "${BLUE}System Logs:${NC}"
    echo -e "${WHITE}=============${NC}"
    echo
    
    echo -e "${CYAN}Recent System Messages (last 15 lines):${NC}"
    if command -v journalctl >/dev/null 2>&1; then
        journalctl -n 15 --no-pager | tail -15
    elif [ -f /var/log/syslog ]; then
        tail -15 /var/log/syslog
    elif [ -f /var/log/messages ]; then
        tail -15 /var/log/messages
    else
        echo "  No accessible system logs found"
    fi
    echo
    
    echo -e "${CYAN}Recent Authentication Logs:${NC}"
    if [ -f /var/log/auth.log ]; then
        tail -5 /var/log/auth.log | grep -v "sudo.*pam_unix"
    elif [ -f /var/log/secure ]; then
        tail -5 /var/log/secure
    else
        echo "  No authentication logs found"
    fi
    
    wait_for_key
}

# Function to handle invalid input
invalid_option() {
    show_header
    echo -e "${RED}Invalid option. Please try again.${NC}"
    sleep 2
}

# Main loop
main_loop() {
    while true; do
        show_menu
        read -n 1 choice
        echo
        
        case $choice in
            1) show_system_info ;;
            2) show_network_status ;;
            3) show_disk_usage ;;
            4) show_services ;;
            5) show_memory_usage ;;
            6) show_processes ;;
            7) show_logs ;;
            0) 
                clear
                echo -e "${GREEN}Thank you for using Ubuntu System Management!${NC}"
                echo -e "${YELLOW}Goodbye!${NC}"
                exit 0
                ;;
            *) invalid_option ;;
        esac
    done
}

# Trap Ctrl+C
trap 'clear; echo -e "\n${YELLOW}Exiting...${NC}"; exit 0' INT

# Check if running as root and warn
if [ "$EUID" -eq 0 ]; then
    echo -e "${YELLOW}Warning: Running as root. Some information may be limited.${NC}"
    sleep 2
fi

# Start the main loop
main_loop
