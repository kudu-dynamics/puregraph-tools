#!/bin/bash

# Load keys from a JSON payload.
read -a arr < <(echo $(jq -c -r '.EventName, .Key, .Missing'))

# Only process events that are "s3:ObjectCreated:*"
if [[ "${arr[0]}" = s3:ObjectCreated:* ]]; then
  echo "eventname: \"${arr[0]}\""
  echo "key: \"${arr[1]}\""
  echo "missing: \"${arr[2]}\""
fi
