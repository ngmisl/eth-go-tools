#!/bin/bash

# Function to prompt the user for a folder name
ask_for_folder_name() {
    read -p "Enter the folder name: " folder_name
}

# Function to create the folder and change into it
create_and_cd_into_folder() {
    mkdir -p "$folder_name"
    cd "$folder_name" || { echo "Failed to change directory to $folder_name"; exit 1; }
}

# Function to initialize a new forge project and install OpenZeppelin contracts
initialize_forge_and_install_openzeppelin() {
    forge init
    forge install OpenZeppelin/openzeppelin-contracts
}

# Function to set remappings in remappings.txt
set_remappings() {
    local remappings_file="remappings.txt"
    local remappings=(
        "@openzeppelin/contracts/=lib/openzeppelin-contracts-upgradeable/lib/openzeppelin-contracts/contracts/"
        "@openzeppelin/contracts-upgradeable/=lib/openzeppelin-contracts-upgradeable/contracts/"
    )

    # Empty the remappings.txt file
    > "$remappings_file"

    # Write the remappings to the file
    for remapping in "${remappings[@]}"; do
        echo "$remapping" >> "$remappings_file"
    done
}

# Main script execution
ask_for_folder_name
create_and_cd_into_folder
initialize_forge_and_install_openzeppelin
set_remappings

echo "Setup complete. You are now in the folder $folder_name and OpenZeppelin contracts are installed with remappings set."
