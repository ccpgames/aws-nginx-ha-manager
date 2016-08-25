package monitor_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/ccpgames/aws-nginx-ha-manager/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigWriter", func() {
	var (
		configWriter *ConfigWriter
		fileFH       *os.File
		ipList       []string
		ipListStr    string
		upstreamName string
	)

	BeforeEach(func() {
		fileFH, err := ioutil.TempFile("", "config_writer_tests")
		if err != nil {
			fmt.Errorf("Error opening temp file", err)
		}
		configWriter = NewConfigWriter(fileFH.Name())
	})

	AfterEach(func() {
		fileFH.Close()
	})

	Describe("Write config", func() {
		Context("Wtih ipList "+ipListStr, func() {
			expected := fmt.Sprintf(`upstream %s {
	%s
}`, upstreamName, strings.Join(ipList, ",\n"))
			It("should write a file containing all the ips", func() {
				configWriter.WriteConfig(ipList)
				Expect(ioutil.ReadFile(fileFH.Name())).To(Equal(expected))
			})
		})
	})
})
