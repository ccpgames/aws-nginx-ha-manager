package monitor

// DbusConnection is a wrapping interface so we can mock dbus.Conn
type DbusConnection interface {
	ReloadOrRestartUnit(name string, mode string, ch chan<- string)
}

// Monitor encapsulates a single runtime, including configurations
type Monitor struct {
	ConfigPath string
	dbusConn   *DbusConnection
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(configPath string, dbusConn *DbusConnection) *Monitor {
	monitor := Monitor{
		ConfigPath: configPath,
		dbusConn:   dbusConn,
	}

	return &monitor
}
