package monitor

// Monitor encapsulates a single runtime, including configurations
type Monitor struct {
	ConfigPath string
}

// NewMonitor return a new instance of Monitor with the supplied parameters
func NewMonitor(configPath string) *Monitor {
	monitor := Monitor{
		ConfigPath: configPath,
	}

	return &monitor
}
