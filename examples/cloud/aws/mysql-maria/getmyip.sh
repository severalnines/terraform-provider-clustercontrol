#!/bin/bash

# Fetch the public IP address and return it in JSON format
IP=$(curl -s https://ifconfig.me)
echo "{\"ip\":\"${IP}\"}"

