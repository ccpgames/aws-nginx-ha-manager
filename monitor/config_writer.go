package monitor

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
)

var templateTxt = `upstream {{.UpstreamName}} {
{{$len := lendec .IPList 1}}
{{range $i, $e := .IPList}}	{{$e}}{{if lt $i $len}},{{end}}
{{end}}}`

// ConfigWriter encapsulates writing an upstream config to a path
type ConfigWriter struct {
	configPath   string
	template     *template.Template
	UpstreamName string
	IPList       []string
}

func lendec(arr []string, n int) int {
	return len(arr) - n
}

func lt(a int, b int) bool {
	return a < b
}

// NewConfigWriter constructs an instance of ConfigWriter
func NewConfigWriter(path string, upstreamName string) *ConfigWriter {
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
	ioutil.WriteFile(w.configPath, buf.Bytes(), 600)
	return nil
}
