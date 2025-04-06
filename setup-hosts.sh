#!/bin/bash

# Define the domain and IP
DOMAIN="github.local"
IP="0.0.0.0"

# Check if the entry already exists
if grep -q "^${IP}\s*${DOMAIN}" /etc/hosts; then
    echo "Host entry already exists"
else
    # Create the entry
    ENTRY="${IP}    ${DOMAIN}"
    
    # Show the command that needs to be run
    echo "Please run the following command to add the host entry:"
    echo "sudo sh -c 'echo \"${ENTRY}\" >> /etc/hosts'"
    
    # Also show how to remove it
    echo -e "\nTo remove the entry later, run:"
    echo "sudo sed -i '/${DOMAIN}/d' /etc/hosts"
fi

# Show current entry in hosts file
echo -e "\nCurrent entries in /etc/hosts for ${DOMAIN}:"
grep "${DOMAIN}" /etc/hosts || echo "No entry found" 