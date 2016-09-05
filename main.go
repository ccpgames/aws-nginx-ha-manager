package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/ccpgames/aws-nginx-ha-manager/cmd"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	host := "vkpilot.cert-valkyrieapi.com"
	addr, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%s: %s", host, addr)
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
