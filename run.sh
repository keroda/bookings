#!/bin/bash

go build -o bookings cmd/web/*.go 
./bookings  -dbname=bookings -dbuser=kjetilrodal -cache=false -production=false