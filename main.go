package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	_ "github.com/tmeckel/coredns-finalizer"
)

func init() {
	dnsserver.Directives = append(
		dnsserver.Directives,
		"finalize",
	)
}

func main() {
	coremain.Run()
}
