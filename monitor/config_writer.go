package monitor

// ConfigWriter encapsulates writing an upstream config to a path
type ConfigWriter struct {
	configPath string
}

// NewConfigWriter constructs an instance of ConfigWriter
func NewConfigWriter(path string) *ConfigWriter {
	cw := ConfigWriter{path}
	return &cw
}

// WriteConfig writes an upstream config to a file based on path
func (w *ConfigWriter) WriteConfig(ipList []string) error {
	return nil
}
