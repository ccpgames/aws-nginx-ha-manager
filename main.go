package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ccpgames/aws-nginx-ha-manager/cmd"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
