package main

import (
	"flag"

	log "github.com/kirillDanshin/dlog"
)

const (
	ver = `4.0 "Dark Dream"`
)

var (
	verHash, verTimeStamp string

	flagWebhook = flag.Bool("webhook", false, "activate getting updates via webhook")
)

func init() {
	log.Ln("Running", ver, "version...")
	flag.Parse()

	cfgInit()
	langInit()
}
