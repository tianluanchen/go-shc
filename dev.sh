#!/usr/bin/env bash

#  You need to install nodejs and nodemon for it to work properly.
nodemon --exec "go run main.go serve -a 127.0.0.1:8080 --trim-path" --ext "go,html,js,css,json" --delay 1s
