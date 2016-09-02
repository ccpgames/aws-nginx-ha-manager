package monitor

import (
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
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(configPath string, dbusConn DbusConnection, interval int, fqdn string) *Monitor {
	resolver := NewAWSResolver()
	monitor := Monitor{
		ConfigPath:   configPath,
		dbusConn:     dbusConn,
		interval:     interval,
		balancer:     NewBalancer(resolver, fqdn),
		configWriter: NewConfigWriter(configPath, "aws_upstream"),
		host:         fqdn,
	}

	return &monitor
}

// IsStopped checks if the monitor is running
func (m *Monitor) IsStopped() bool {
	return m.stop
}

// Loop runs the monitor and resolves the service at a given interval
func (m *Monitor) Loop(ch chan syscall.Signal) {
	var ipList []string
	var sleepyTime time.Duration
	sleepyTime = time.Duration(m.interval) * time.Millisecond
	for !m.stop {
		// Lets get the ip list
		list, err := m.balancer.GetIPList()
		if err != nil {
			log.Errorf("Error getting ip list: %s", err)
		}
		// Do we have new list
		if !testEq(ipList, list) {
			ipList = list
			m.configWriter.WriteConfig(ipList)
			retVal, err := m.dbusConn.ReloadOrRestartUnit("nginx", "fail", nil)
			if err != nil {
				log.Errorf("Error restarting or reloading the nginx unit: retval: %d: %s", retVal, err)
			}
			ch <- syscall.Signal(0)
		}
		// Check for signals from outside
		select {
		case s := <-ch:
			switch s {
			case syscall.SIGHUP:
				retVal, err := m.dbusConn.ReloadOrRestartUnit("nginx", "fail", nil)
				if err != nil {
					log.Errorf("Error restarting or reloading the nginx unit: retval: %d: %s", retVal, err)
				}
				log.Info("Reloaded nginx by signal")
				ch <- syscall.SIGHUP
				break
			case syscall.SIGKILL:
			case syscall.SIGABRT:
				log.Info("Exiting by request")
				m.stop = true
				ch <- syscall.SIGABRT
			}
		default:
		}
		// Wait for a bit
		if !m.stop {
			time.Sleep(sleepyTime)
		}
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
