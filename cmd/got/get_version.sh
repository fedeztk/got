#!/usr/bin/env bash

printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)" > .version
