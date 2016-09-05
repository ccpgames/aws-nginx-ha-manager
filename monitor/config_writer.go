package monitor

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

var templateTxt = `upstream {{.UpstreamName}} {
{{$port := .Port}}{{range $i, $e := .IPList}}    server {{$e}}:{{$port}};
{{end}}}`

// ConfigWriter encapsulates writing an upstream config to a path
type ConfigWriter struct {
	configPath   string
	template     *template.Template
	UpstreamName string
	IPList       []string
	Port         int
}

func lendec(arr []string, n int) int {
	return len(arr) - n
}

func lt(a int, b int) bool {
	return a < b
}

// NewConfigWriter constructs an instance of ConfigWriter
func NewConfigWriter(path string, upstreamName string, port int) *ConfigWriter {
	funcMap := template.FuncMap{"lendec": lendec}

	tpl, err := template.New("upstream").Funcs(funcMap).Parse(templateTxt)

	if err != nil {
		log.Fatal(err)
	}
	cw := ConfigWriter{
		configPath:   path,
		template:     tpl,
		UpstreamName: "aws_upstream",
		IPList:       []string{},
		Port:         port,
	}
	return &cw
}

// WriteConfig writes an upstream config to a file based on path
func (w *ConfigWriter) WriteConfig(IPList []string) error {
	w.IPList = IPList
	var buf bytes.Buffer
	err := w.template.Execute(&buf, w)
	if err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	err = ioutil.WriteFile(w.configPath, buf.Bytes(), os.FileMode(0644))
	if err != nil {
		log.Errorf("Error writing config file: %s", err)
		return err
	}
	log.Infof("Wrote new upstream config to %s", w.configPath)
	return nil
}
