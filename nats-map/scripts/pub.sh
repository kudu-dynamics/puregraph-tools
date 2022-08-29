#!/bin/sh

. ./scripts/test.env

FLIP=$((RANDOM%2))
if [ $FLIP -eq 0 ]; then
  PAYLOAD='{"EventName": "s3:ObjectCreated:Put", "Key": "bucket/key"}'
else
  PAYLOAD='{"EventName": "s3:ObjectAccessed:Get", "Key": "bucket/key"}'
fi

stan-pub \
  -s "${NATS_URL}" \
  -c "${STAN_CLUSTER}" \
  -id "${STAN_CLIENT_ID}-pub-1" \
  "${STAN_CHANNEL}" "${PAYLOAD}"
