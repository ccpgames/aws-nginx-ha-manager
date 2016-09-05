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
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/ccpgames/aws-nginx-ha-manager/monitor"
	"github.com/coreos/go-systemd/dbus"
	"github.com/spf13/cobra"
)

var elbName string
var upstreamName string
var port int
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
		elbName = args[0]
		if dbusConn, err = dbus.New(); err != nil {
			log.Fatalf("Could not open dbus connection (this program requires linux with nginx runnin on systemd): %s", err)
		}
		if _, err = dbusConn.GetUnitProperties("nginx.service"); err != nil {

			units, err := dbusConn.ListUnits()
			unitNames := make([]string, len(units))
			for i, _ := range units {
				unitNames[i] = units[i].Name
			}
			if err != nil {
				log.Fatalf("Error getting unitlist: %s", err)
			}
			log.Fatalf("Could not get properties of nginx unit; is it running?: %s (available units listed below)\n%s", err, strings.Join(unitNames, "\n"))
		}
		resolver := monitor.NewAWSResolver()
		mon := monitor.NewMonitor(configFile, dbusConn, interval, elbName, port, upstreamName, resolver)
		ch := make(chan syscall.Signal)
		mon.Loop(ch)
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
	monitorCmd.PersistentFlags().IntVar(&interval, "interval", 5, "the interval in seconds to poll the elbName for new upstream hosts")
	monitorCmd.PersistentFlags().StringVar(&configFile, "upstream-file", "/etc/nginx/conf.d/aws_upstream.conf", "The upstream config file to write to")
	monitorCmd.PersistentFlags().IntVar(&port, "port", 10080, "The port upstream servers are called on")
	monitorCmd.PersistentFlags().StringVar(&upstreamName, "upstream-name", "upstream", "the name of the upstream")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
