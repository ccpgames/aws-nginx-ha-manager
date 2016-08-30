package monitor_test

import (
	"fmt"
	"io/ioutil"
	"log"
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
		err          error
		expected     string
	)

	BeforeEach(func() {
		fileFH, err = ioutil.TempFile("", "config_writer_tests")
		if err != nil {
			fmt.Errorf("Error opening temp file", err)
		}
		configWriter = NewConfigWriter(fileFH.Name(), upstreamName)
		log.Printf("using file %s", fileFH.Name())
		ipList = []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"}
		ipListStr = strings.Join(ipList, ", ")
		upstreamName = "aws_upstream"
		expected = fmt.Sprintf(`upstream %s {
	%s
}`, upstreamName, strings.Join(ipList, ",\n\t"))
	})

	AfterEach(func() {
		fileFH.Close()
	})

	Describe("Write config", func() {
		Context("With ipList", func() {
			It("should write a file containing all the ips", func() {
				configWriter.WriteConfig(ipList)
				actual, err := ioutil.ReadFile(fileFH.Name())
				if err != nil {
					log.Fatalf("Error reading file %s: %s", fileFH.Name(), err)
				}
				Expect(string(actual)).To(Equal(expected))
			})
		})
	})
})
