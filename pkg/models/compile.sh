#!/bin/sh
easyjson -all export.go
easyjson -all import.go
easyjson -all rpl.go
easyjson -all title.go
