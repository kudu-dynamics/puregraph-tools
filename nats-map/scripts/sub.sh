#!/bin/sh

. ./scripts/test.env

export STAN_CLIENT_ID="${STAN_CLIENT_ID}-sub-1"

./nats-map
