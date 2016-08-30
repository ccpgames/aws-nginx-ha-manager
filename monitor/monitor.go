package monitor

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
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(configPath string, dbusConn *DbusConnection, interval int, fqdn string) *Monitor {
	monitor := Monitor{
		ConfigPath:   configPath,
		dbusConn:     dbusConn,
		interval:     interval,
		balancer:     NewBalancer(fqdn),
		configWriter: NewConfigWriter(configPath, "aws_upstream"),
	}

	return &monitor
}

// Loop runs the monitor and resolves the service at a given interval
func (m *Monitor) Loop(interval int) <-chan bool {
	quit := make(chan bool)

	return quit
}
