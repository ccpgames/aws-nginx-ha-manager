package monitor

import (
	"log"
	"time"
)

// DbusConnection is a wrapping interface so we can mock dbus.Conn
type DbusConnection interface {
	ReloadOrRestartUnit(name string, mode string, ch chan<- string)
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
func NewMonitor(configPath string, dbusConn *DbusConnection, interval int, fqdn string) *Monitor {
	resolver := NewAWSResolver()
	monitor := Monitor{
		ConfigPath:   configPath,
		dbusConn:     *dbusConn,
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
	var ipList []string
	var sleepyTime time.Duration
	sleepyTime = time.Duration(interval) * time.Second
	for !m.stop {
		list, err := m.balancer.GetIPList()
		if err != nil {
			log.Printf("Error getting ip list: %s", err)
		}
		if !testEq(ipList, list) {
			ipList = list
			m.configWriter.WriteConfig(ipList)
			m.dbusConn.ReloadOrRestartUnit("nginx", "fail", nil)
		}
		time.Sleep(sleepyTime)
	}
	return nil
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
