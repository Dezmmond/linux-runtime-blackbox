#!/usr/bin/env bash
set -euo pipefail

sleep 2 &
child_one=$!

sh -c 'sleep 2' &
child_two=$!

wait "$child_one"
wait "$child_two"
