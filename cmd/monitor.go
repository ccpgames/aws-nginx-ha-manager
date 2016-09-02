// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/ccpgames/aws-nginx-ha-manager/monitor"
	"github.com/coreos/go-systemd/dbus"
	"github.com/spf13/cobra"
)

var fqdn string
var interval int
var configFile string

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring an LB",
	Long:  `nginx-aws-monitor monitors a Load Balancer and updates an nginx upstream configuration file to match the results`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		var err error
		var dbusConn *dbus.Conn
		if dbusConn, err = dbus.New(); err != nil {
			log.Fatalf("Could not open dbus connection (this program requires linux with nginx runnin on systemd): %s", err)
		}
		monitor := monitor.NewMonitor(configFile, dbusConn, interval, fqdn)
		ch := make(chan syscall.Signal)
		monitor.Loop(ch)
		run := true
		for run {
			select {
			case sig := <-ch:
				switch sig {
				case syscall.SIGHUP:
					log.Info("Reloaded on signal")
					break
				case syscall.SIGABRT:
					log.Info("Exiting on signal")
					run = false
					break
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")
	monitorCmd.PersistentFlags().IntVar(&interval, "interval", 5, "the interval in seconds to poll the fqdn for new upstream hosts")
	monitorCmd.PersistentFlags().StringVar(&fqdn, "fqdn", "", "Specify the fqdn to monitor")
	monitorCmd.PersistentFlags().StringVar(&configFile, "upstream-file", "/etc/nginx/conf.d/aws_upstream.conf", "The upstream config file to write to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
