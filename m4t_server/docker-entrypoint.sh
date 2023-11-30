#!/usr/bin/env bash
#
#
set -e


# if command starts with an option, prepend omni-server
if [ "${1:0:1}" = '-' ]; then
    set -- ./serve.py "$@"
fi
# cd workspace


# if command app only, add use default args
if [ "$1" = './serve.py' ] && [ "$#" -eq 1 ]; then
    exec ./serve.py  --host "0.0.0.0"  --port ${SERVER_PORT} --model-path ${MODEL_PATH} 
fi

exec "$@"
