#!/bin/bash

# kill server on shutdown
trap 'pkill -P $$' 0

start_server() {
	make run&
	pid=$!
}

restart_server() {
	kill $pid
	start_server
}

start_server

while true; do
	# restart if any .go files are modified
	if [ `inotifywait --event create,modify,delete --recursive --format '%w%f' . 2> /dev/null | grep -E '\.go$'` ]; then
		restart_server
	fi
done
