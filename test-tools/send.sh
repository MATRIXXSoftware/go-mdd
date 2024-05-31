#!/bin/bash

server_address="127.0.0.1"
server_port="14060"

hex_string=$(<msg.hex)

binary_data=$(echo "$hex_string" | xxd -r -p)

echo -n "$binary_data" | nc "$server_address" "$server_port"

