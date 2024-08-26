package main

// go build -ldflags="-s -w -X ginapp/conf.BuildDate=$(date -Iseconds) -o 2wayssl .
var (
	BuildCommitId = "000"
	BuildDate     = ""
	GoVersion   = ""
	BuildVersion  = ""
)
