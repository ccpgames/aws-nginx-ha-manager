package monitor

import (
	"os"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

// DbusConnection is a wrapping interface so we can mock dbus.Conn
type DbusConnection interface {
	ReloadOrRestartUnit(name string, mode string, ch chan<- string) (int, error)
}

// Monitor encapsulates a single runtime, including configurations
type Monitor struct {
	ConfigPath   string
	dbusConn     DbusConnection
	interval     int
	configWriter *ConfigWriter
	balancer     *Balancer
	stop         bool
	host         string
	port         int
	upstreamName string
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(
	configPath string,
	dbusConn DbusConnection,
	interval int,
	elbName string,
	port int,
	upstreamName string,
	resolver Resolver,
) *Monitor {
	monitor := Monitor{
		ConfigPath:   configPath,
		dbusConn:     dbusConn,
		interval:     interval,
		balancer:     NewBalancer(resolver, elbName),
		configWriter: NewConfigWriter(configPath, upstreamName, port),
		host:         elbName,
		port:         port,
		upstreamName: upstreamName,
	}

	return &monitor
}

// IsStopped checks if the monitor is running
func (m *Monitor) IsStopped() bool {
	return m.stop
}

// Loop runs the monitor and resolves the service at a given interval
func (m *Monitor) Loop(sig chan os.Signal, msg chan string) {
	var ipList []string
	var sleepyTime time.Duration
	sleepyTime = time.Duration(m.interval) * time.Millisecond * 1000
	for !m.stop {
		// Lets get the ip list
		list, err := m.balancer.GetIPList()
		if err != nil {
			log.Errorf("Error getting ip list: %s", err)
		}
		// Do we have new list
		if !testEq(ipList, list) {
			log.Infoln("IP List updated, writing new configuration")
			ipList = list
			m.configWriter.WriteConfig(ipList)
			retVal, err := m.dbusConn.ReloadOrRestartUnit("nginx.service", "fail", nil)
			if err != nil {
				log.Errorf("Error restarting or reloading the nginx unit: retval: %d: %s", retVal, err)
			}
			msg <- "Updated and reloaded configuration"
		} else {
			log.Debugln("No changes detected")
		}
		// Check for signals from outside
		select {
		case s := <-sig:
			switch s {
			case syscall.SIGHUP:
				retVal, err := m.dbusConn.ReloadOrRestartUnit("nginx.service", "fail", nil)
				if err != nil {
					log.Errorf("Error restarting or reloading the nginx unit: retval: %d: %s", retVal, err)
				}
				log.Info("Reloaded nginx by signal")
				msg <- "Reloaded configuration"
				break
			case os.Interrupt:
				log.Info("Exiting by request")
				m.stop = true
				msg <- "Exit"
			}
		default:
			log.Debugln("No signal")
		}
		if m.stop {
			// Break immediately on signal
			break
		}
		// Wait for a bit
		log.Debugf("Sleeping for %v seconds", m.interval)
		time.Sleep(sleepyTime)
	}
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
