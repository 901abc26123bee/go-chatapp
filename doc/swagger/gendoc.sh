#!/bin/bash

# Function to initialize swag in a given directory with a specific controller
initialize_swag() {
  local dir=$1
  local controller=$2
  
  # Navigate to the directory
  cd "$dir" || exit

  # Initialize swag with the specified controller
  swag init -g "$controller"

  # Return to the original directory
  cd - > /dev/null || exit
}

# TODO: refine input
# Initialize swag for the account directory
initialize_swag "./doc/swagger/account" "../../../api/account/controller.go"

# Initialize swag for the realtime directory
initialize_swag "./doc/swagger/realtime" "../../../api/realtime/controller.go"
