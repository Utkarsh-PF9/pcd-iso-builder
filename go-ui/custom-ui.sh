#!/bin/bash
# Run system update
sudo apt update

# ⚠️ Run the UI binary script (ensure path is correct and executable)
if [[ -x /usr/local/bin/pcd-iso-ui ]]; then
    /usr/local/bin/pcd-iso-ui
else
    echo "Error: /usr/local/bin/pcd-iso-ui not found or not executable."
    exit 1
fi

# After exit, show message
echo "Custom UI has exited."
echo "You now have access to the normal Ubuntu terminal."
echo "Type 'custom-ui' to restart the management interface."

exit 0
