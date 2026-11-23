#!/bin/bash
set -u

txshield-api -serve -addr 127.0.0.1:8080 &
backend_pid=$!
node ./node_modules/next/dist/bin/next start --hostname 0.0.0.0 --port 3000 &
frontend_pid=$!

stop_servers() {
    kill -TERM "$backend_pid" "$frontend_pid" 2>/dev/null || true
    wait "$backend_pid" "$frontend_pid" 2>/dev/null || true
}

trap 'stop_servers; exit 0' TERM INT
wait -n "$backend_pid" "$frontend_pid"
status=$?
stop_servers
exit "$status"
