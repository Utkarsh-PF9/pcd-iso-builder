#!/bin/bash

# Check if Python3 is available
if ! command -v python3 &> /dev/null; then
    echo "Python3 is required but not installed."
    echo "Installing Python3..."
    sudo apt update && sudo apt install -y python3
fi

# Run the Python curses interface
python3 /usr/local/bin/custom-ui.py

# After Python script exits, we're back in the shell
echo "Custom UI has exited."
echo "You now have access to the normal Ubuntu terminal."
echo "Type 'custom-ui' to restart the management interface."

exit 0
