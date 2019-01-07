#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd "${SCRIPTPATH}"/../
find githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
