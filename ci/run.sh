#!/usr/bin/env bash
set -x
set -e

CONFIG=$1 ; shift

function finish() {
	set +e
	kill %2
	kill %1
	wait
}

$PWD/hege-front \
	-templates $PWD/front-server/templates \
	-static $PWD/front-server/static \
	&

$PWD/hege-world \
	-load $CONFIG \
	&

trap finish SIGTERM SIGINT
wait
