#!/bin/bash

# Clean up .custos data for all users
echo "Cleaning up Custos data..."

# Remove from all user home directories
for user_home in /home/*; do
    if [ -d "$user_home/.custos" ]; then
        echo "Removing data for user home: $user_home"
        rm -rf "$user_home/.custos"
    fi
done

# Also check root
if [ -d "/root/.custos" ]; then
    echo "Removing data for root"
    rm -rf "/root/.custos"
fi

exit 0
