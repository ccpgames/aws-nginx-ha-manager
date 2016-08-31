package monitor

import (
	"net"
)

// DbusConnection is a wrapping interface so we can mock dbus.Conn
type DbusConnection interface {
	ReloadOrRestartUnit(name string, mode string, ch chan<- string)
}

// Monitor encapsulates a single runtime, including configurations
type Monitor struct {
	ConfigPath   string
	dbusConn     *DbusConnection
	interval     int
	configWriter *ConfigWriter
	balancer     *Balancer
	stop         bool
	host         string
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(configPath string, dbusConn *DbusConnection, interval int, fqdn string) *Monitor {
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

// Break exits Loop cleanly
func (m *Monitor) Break() {
	m.stop = true
}

// Loop runs the monitor and resolves the service at a given interval
func (m *Monitor) Loop(interval int) (err error) {
	for !m.stop {
		net.LookupHost(m.host)
	}
	return nil
}
