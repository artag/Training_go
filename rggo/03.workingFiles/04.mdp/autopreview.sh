#! /bin/bash

# This script receives the name of the file you want to preview as an argument.
# It calculates the checksum of this file every 5 seconds.
# If the file was changed, run mdp tool to preview it.

# Make the script executable:
# chmod +x autopreview.sh

# Run script example:
# ./autopreview.sh file.md
#
# Stop script:
# CTRL + C

FHASH=`md5sum $1`
while true; do
    NHASH=`md5sum $1`
    if [ "$NHASH" != "$FHASH" ]; then
        ./mdp -file $1
        FHASH=$NHASH
    fi
    sleep 5
done
