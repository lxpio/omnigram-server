#!/usr/bin/env sh
#
#
set -e

# if command starts with an option, prepend omni-server
if [ "${1:0:1}" = '-' ]; then
    set -- omni-server "$@"
fi
# cd workspace

# if command app only, add use default args
if [ "$1" = 'omni-server' ] && [ "$#" -eq 1 ]; then
    exec omni-server -conf ${CONFIG_FILE}
fi

exec "$@"
